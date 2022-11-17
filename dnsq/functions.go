package dnsq

import (
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/miekg/dns"
	"github.com/rodaine/table"
)

func CheckSOA(domain string, nameserver string) uint32 {

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{fmt.Sprintf("%v.", domain), dns.TypeSOA, dns.ClassINET}
	c := new(dns.Client)
	in, _, err := c.Exchange(m1, fmt.Sprintf("%v:53", nameserver))

	if err == nil {
		if t, ok := in.Answer[0].(*dns.SOA); ok {
			//fmt.Printf("%v: %v\n", nameserver, t.Serial)
			return t.Serial
		} else {
			return 0
		}
	} else {
		return 0
	}

}

func GetNS(domain string, nameserver string) []string {

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{fmt.Sprintf("%v.", domain), dns.TypeNS, dns.ClassINET}
	c := new(dns.Client)
	in, _, err := c.Exchange(m1, fmt.Sprintf("%v:53", nameserver))

	var nameservers []string

	if err == nil {
		for _, x := range in.Answer {
			if t, ok := x.(*dns.NS); ok {
				nameservers = append(nameservers, t.Ns)
			}
		}

	}
	return nameservers
}

func Runner(domain string, ns string, nameservers []string) {
	var results []Result

	for _, nameserver := range nameservers {
		serial := CheckSOA(domain, nameserver)

		results = append(results, Result{
			Domain:     domain,
			Nameserver: nameserver,
			Serial:     serial,
		})

	}

	var highest_serial uint32 = 0

	// Determine highest serial
	for _, r := range results {
		if r.Serial > highest_serial {
			highest_serial = r.Serial
		}
	}

	// Calculate difference
	for k, v := range results {

		results[k].Diff = highest_serial - v.Serial
	}

	// Nice table output
	now := time.Now()

	fmt.Printf("Run executed at: %v\n\n", now)
	headerFmt := color.FgGreen.Sprintf

	tbl := table.New("Domain", "Nameserver", "Serial", "Difference")

	tbl.WithHeaderFormatter(headerFmt)

	green := color.FgGreen.Render
	yellow := color.FgYellow.Render
	red := color.FgRed.Render

	for _, line := range results {
		// GREEN
		if line.Diff >= 0 && line.Diff <= 100 {
			tbl.AddRow(line.Domain, line.Nameserver, line.Serial, green(line.Diff))
		}
		// YELLOW
		if line.Diff >= 101 && line.Diff <= 200 {
			tbl.AddRow(line.Domain, line.Nameserver, line.Serial, yellow(line.Diff))
		}
		// RED
		if line.Diff >= 201 {
			tbl.AddRow(line.Domain, line.Nameserver, line.Serial, red(line.Diff))
		}
	}

	tbl.Print()
}
