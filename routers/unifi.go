package routers

type UnifiConfig struct {
}

type UnifiReadPortForward struct {
	PFWDInterface   string `json:"pfwd_interface"` // Port forward interface. Options: wan, wan2, both
	Forward         string `json:"fwd"`            // IP to forward from
	DestinationIP   string `json:"destination_ip"` // Destination IP
	Source          string `json:"src"`            // Valid source ips "any" for any
	Log             bool   `json:"log"`            // Enable logging
	Protocol        string `json:"proto"`          //Protocol tcp, udp, tcp_udp
	Name            string `json:"name"`           // Friendly name
	DestinationPort string `json:"dst_port"`       // Destination Port
	SiteID          string `json:"site_id"`        //Site ID of the controller
	ID              string `json:"_id"`            //ID of the port forward?
	ForwardPort     string `json:"fwd_port"`       //Forward port
	Enabled         bool   `json:"enabled"`        //Enabled
}
