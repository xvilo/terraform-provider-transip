package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/transip/gotransip"
	"github.com/transip/gotransip/domain"
)

func resourceDNSEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSEntryCreate,
		Read:   resourceDNSEntryRead,
		Update: resourceDNSEntryUpdate,
		Delete: resourceDNSEntryDelete,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// TODO: true for transip?
				StateFunc: func(v interface{}) string {
					value := strings.TrimSuffix(v.(string), ".")
					return strings.ToLower(value)
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expire": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  86400,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(domain.DNSEntryTypeA),
					string(domain.DNSEntryTypeAAAA),
					string(domain.DNSEntryTypeCNAME),
					string(domain.DNSEntryTypeMX),
					string(domain.DNSEntryTypeNS),
					string(domain.DNSEntryTypeTXT),
					string(domain.DNSEntryTypeSRV),
				}, false),
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDNSEntryCreate(d *schema.ResourceData, m interface{}) error {
	domainName := d.Get("domain")

	entry := domain.DNSEntry{
		d.Get("name").(string),
		d.Get("expire").(int64),
		d.Get("type").(domain.DNSEntryType),
		d.Get("content").(string),
	}

	id := fmt.Sprintf("%s-%s-%s-%s-%s",
		domainName, entry.Name, entry.TTL, entry.Type, entry.Content)
	d.SetId(id)

	return resourceDNSEntryRead(d, m)
}

func resourceDNSEntryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	domainName := d.Get("domain").(string)
	name := d.Get("name").(string)

	dom, err := domain.GetInfo(client, domainName)
	if err != nil {
		return fmt.Errorf("failed to get domain %s for reading entry: %s", domainName, err)
	}

	entries := []domain.DNSEntry{}
	for _, e := range dom.DNSEntries {
		if e.Name == name {
			entries = append(entries, e)
		}
	}
	if len(entries) == 0 {
		d.SetId("")
		return nil
	}
	if len(entries) > 1 {
		return fmt.Errorf("multirecord not yet supported")
	}
	entry := entries[0]

	d.Set("name", entry.Name)
	d.Set("expire", entry.TTL)
	d.Set("type", entry.Type)
	d.Set("content", entry.Content)
	return nil
}

func resourceDNSEntryUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDNSEntryRead(d, m)
}

func resourceDNSEntryDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
