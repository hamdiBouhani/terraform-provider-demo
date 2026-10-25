package provider

import (
	"context"
	"strconv"

	"github.com/hamdiBouhani/terraform-provider-demo/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DEMO_BASE_URL", "https://jsonplaceholder.typicode.com"),
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DEMO_API_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"demo_post": resourcePost(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	c := client.New(d.Get("base_url").(string), d.Get("api_key").(string))
	return c, nil
}

func resourcePost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePostCreate,
		ReadContext:   resourcePostRead,
		UpdateContext: resourcePostUpdate,
		DeleteContext: resourcePostDelete,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"body": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	post, err := c.CreatePost(ctx, d.Get("title").(string), d.Get("body").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(post.ID))
	return resourcePostRead(ctx, d, meta)
}

func resourcePostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id, _ := strconv.Atoi(d.Id())
	post, err := c.GetPost(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if post == nil {
		d.SetId("") // removed outside terraform
		return nil
	}

	d.Set("title", post.Title)
	d.Set("body", post.Body)
	return nil
}

func resourcePostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// For simplicity, delete + recreate (or implement PUT)
	if d.HasChanges("title", "body") {
		c := meta.(*client.Client)
		_ = c.DeletePost(ctx, mustAtoi(d.Id()))
		post, err := c.CreatePost(ctx, d.Get("title").(string), d.Get("body").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(strconv.Itoa(post.ID))
	}
	return resourcePostRead(ctx, d, meta)
}

func resourcePostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	if err := c.DeletePost(ctx, mustAtoi(d.Id())); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
