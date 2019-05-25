package azurecompute

type imageOS struct {
	Publisher string
	Offer     string
	Sku       string
	Version   string
}

func Image(image string) imageOS {

	var p imageOS

	switch image {
	case "ubuntu":
		p = imageOS{Publisher: "Canonical", Offer: "UbuntuServer", Sku: "16.04-LTS", Version: "latest"}
	case "centos":
		p = imageOS{Publisher: "OpenLogic", Offer: "CentOS", Sku: "7.1", Version: "latest"}
	case "rhel":
		p = imageOS{Publisher: "RedHat", Offer: "RHEL", Sku: "7.2", Version: "latest"}
	case "windows":
		p = imageOS{Publisher: "MicrosoftWindowsServer", Offer: "WindowsServer", Sku: "2016-Datacenter", Version: "latest"}
	}
	return p
}
