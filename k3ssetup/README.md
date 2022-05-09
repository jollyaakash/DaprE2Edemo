## SETUP 

#### Create SSH Keys
```
ssh-keygen -m PEM -t rsa -b 4096 -f ubuntuvm
mv ubuntuvm ubuntuvm.pub ~/.ssh/
eval `ssh-agent -s`
ssh-add ~/.ssh/ubuntuvm
```

### Create Main node with empty vnet
```
az vm create --resource-group e2edemo --name mainnode --image UbuntuLTS --admin-usern
ame aakashm --ssh-key-values ~/.ssh/ubuntuvm.pub --location eastus --size Standard_D4as_v4 --public-ip-sku Standar
d --nsg-rule **none** --vnet-name **e2edemovnet** [ Vnet name can be anything ]
```

#### Get VNET id for other VMs
az network vnet subnet show -g e2edemo --name mainnodeSubnet --vnet-name e2edemovnet --query id -o tsv

### Create Agent node 1 and append VNET id at the end

```
az vm create --resource-group e2edemo --name agentnode1 --image UbuntuLTS --admin-username aakasha1 --ssh-key-values ~/.ssh/ubuntuvm.pub --location eastus --size Standard_D4as_v4 --public-ip-sku Standard --nsg-rule none --subnet /subscriptions/249c0d61-388d-45ad-ba35-f9899f4c1374/resourceGroups/e2edemo/providers/Microsoft.Network/virtualNetworks/e2edemovnet/subnets/mainnodeSubnet

az vm create --resource-group e2edemo --name agentnode2 --image UbuntuLTS --admin-username aakasha2 --ssh-key-values ~/.ssh/ubuntuvm.pub --location eastus --size Standard_D4as_v4 --public-ip-sku Standard --nsg-rule none --subnet /subscriptions/249c0d61-388d-45ad-ba35-f9899f4c1374/resourceGroups/e2edemo/providers/Microsoft.Network/virtualNetworks/e2edemovnet/subnets/mainnodeSubnet
```

### PUBLIC IP
```
mainnode - 20.121.254.241 - aakashm
agentnode1 - 20.232.85.191 - aakasha1
agentnode2 - 20.84.27.100 - aakasha2
clinode - 20.232.52.168 - aakash
```

### vim /etc/hosts
```
10.0.0.4 main - aakashm
10.0.0.5 agent01 - aakasha1
10.0.0.6 agent02 - aakasha1
10.0.0.7 clinode - aakash
```

### Script on every node

sudo apt-get update && sudo apt-get upgrade -y && sudo systemctl reboot

#### Install Docker
```
sudo apt-get remove docker docker-engine docker.io containerd runc && sudo apt-get update && sudo apt-get install ca-certificates curl gnupg lsb-release -y
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update && sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin -y
sudo chmod 666 /var/run/docker.sock
sudo systemctl start docker && sudo systemctl enable docker
```

#### Main node 
```
curl -sfL https://get.k3s.io | sh -s - server --docker --tls-san 10.0.0.7 [ --tls-san flag is used for adding a new hostname or IP for talking to api server]
sudo cat /var/lib/rancher/k3s/server/node-token
```

Token - K100984e007d4c8be843ec9545492d3af72337ea845edd9affd7fea0967aa4a6975::server:07834ee2ec65116a913a836b1e10f2b5

#### Agent node
```
curl -sfL https://get.k3s.io | K3S_URL=https://main:6443 K3S_TOKEN=K100984e007d4c8be843ec9545492d3af72337ea845edd9affd7fea0967aa4a6975::server:07834ee2ec65116a913a836b1e10f2b5 sh -s - --docker

sudo systemctl status k3s-agent
```

## SSH 
ssh -i ~/.ssh/ubuntuvm aakash@ip_address 

## K3s file path 
/etc/rancher/k3s/k3s.yaml
```
scp -i ~/.ssh/ubuntuvm aakashm@20.121.254.241:/etc/rancher/k3s/k3s.yaml dest_file_path
```

## Uninstall 
sudo /usr/local/bin/k3s-agent-uninstall.sh
sudo rm -rf /var/lib/rancher/

### Install Extension 
```
az k8s-extension create -g aajolly --cluster-name e2eDemoArc --cluster-type connectedClusters --name E4K --extension-type microsoft.azedge.mqtt --scope cluster --release-train dev --version 0.1.0
```

## HELP pages

https://computingforgeeks.com/install-kubernetes-on-ubuntu-using-k3s/
