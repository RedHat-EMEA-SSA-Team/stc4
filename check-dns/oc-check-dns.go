package main

import (
    "net"
    "fmt"
    "os"
    "strings"
    flag "github.com/spf13/pflag"
)

func srv(domain string){
    cname, srvs, err := net.LookupSRV("etcd-server-ssl", "tcp", domain)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not get SRV: %v\n", err)
        
    }
    fmt.Printf("\nChecking SRV Records: %s\n\n", cname)

    for _, srv := range srvs {
        if  ( srv.Port == 2380 && srv.Priority == 0 && srv.Weight == 10 && strings.Contains(srv.Target, "etcd-")) && strings.Contains(srv.Target, domain){
        ips, err := net.LookupIP(srv.Target)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
                
            }

            for _, ip := range ips {
                addr, err := net.LookupAddr(ip.String())
                if err != nil {
                        fmt.Fprintf(os.Stderr, "Could not get the fqdn: %v\n", err)
                        
                }

                for i, _ := range addr {
                        if ( len(addr) > 1 ) {
                                if (! strings.Contains(addr[i],srv.Target)) {
                                        continue
                                } else {
                                        fmt.Printf("%s FAIL\n", srv.Target)
                                }
                        } else {
                                if (strings.Contains(addr[0],srv.Target)) {
                                        fmt.Printf("%s FAIL\n", srv.Target)
                                        break
                                } else {
                                        fmt.Printf("%s OK\n", srv.Target)
                                }}
                }
            }

        }
    }
}

func check_nodes(server string,domain string) {
    fmt.Printf("Checking A and PTR for node:...")
    ips, err := net.LookupIP(server+"."+domain)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
        
    }

    for _, ip := range ips {
        addr, err := net.LookupAddr(ip.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "Could not get the fqdn: %v\n", err)
            
        }
        if (strings.Contains(addr[0],server+"."+domain)) {
            fmt.Printf("%s OK\n", server)
        } else {
            fmt.Printf("%s FAIL\n", server)
        }
    }
}

func check_api(domain string) {
    fmt.Printf("\nChecking API and API-INT:...\n\n")
    api, err := net.LookupIP("api."+domain)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
        
    }else {
        fmt.Printf("api.%s - %s: OK\n", domain,api[0])
    }
    api_int, err := net.LookupIP("api-int."+domain)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
        
    }else {
        fmt.Printf("api-int.%s - %s: OK\n", domain,api_int[0])
    }
}

func check_apps(domain string) {
    fmt.Printf("\nChecking *.APPS:...\n\n")
    apps_01, err := net.LookupIP("qwertasdf.apps."+domain)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
        
    }
    apps_02, err := net.LookupIP("poiuylkjh.apps."+domain)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
        
    } 
    if (strings.Contains(apps_01[0].String(),apps_02[0].String())) {
        fmt.Printf("*.apps.%s: OK\n", domain)
    }

}

func main() {
    fmt.Printf("\n######################################################\n")
    fmt.Printf("###     OCP EMEA TOOLKIT - DNS Checking            ###")
    fmt.Printf("\n######################################################\n\n")


    domainPtr := flag.String("domain","","--domain=<domain >")
    nodePtr := flag.String("nodes","","--nodes=master001,master002,master003,worker001,...")
    etcPtr := flag.Bool("etcd",false,"--etcd=true")
    apiPtr := flag.Bool("api",false,"--api=true")
    appsPtr := flag.Bool("apps",false,"--apps=true")
    help := flag.Bool("help",false,"--help Show this!")
  
    flag.Parse()

    if *help {
        flag.Usage()
        return 
    }

    if len(*domainPtr) <= 0 {
        flag.Usage()
        os.Exit(1)
    }

    if len(*domainPtr) > 0 && len(*nodePtr) > 0 {
        fmt.Printf("Domain: %s\n\n", *domainPtr)
        if strings.Contains(*nodePtr,","){
            nodes := strings.Split(*nodePtr, ",")
                for i := range nodes {
                    check_nodes(nodes[i],*domainPtr)
                }
        }else{
            check_nodes(*nodePtr,*domainPtr)
        }
    }
    
    if *etcPtr { 
        srv(*domainPtr)
    }
    
    if *apiPtr {
       check_api(*domainPtr)
    }

    if *appsPtr {
        check_apps(*domainPtr)
     }

     fmt.Printf("\n######################################################\n\n")

}
