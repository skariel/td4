cloud1
root@144.202.104.6
passwd: ****

skariel
skariel@144.202.104.6
passwd: ****

adduser skariel
usermod -aG sudo skariel

simple http:
mkdir pub
cd pub
python3 -m http.server PORT

sudo ufw enable
sudo ufw allow 80/tcp
sudo ufw allow https
sudo ufw allow ssh
sudo ufw status
https://help.ubuntu.com/community/UFW

sudo apt install zsh
oh my zsh:
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"


certbot:
sudo apt-get install certbot
sudo certbot certonly --standalone --dry-run -d solvemytest.dev -d www.solvemytest.dev -d api.solvemytest.dev
cert in: /etc/letsencrypt/live/solvemytest.dev/fullchain.pem
key in: /etc/letsencrypt/live/solvemytest.dev/privkey.pem
renew: certbot renew

vultr nameservers:
----------------------

(insert in customDNS in namecheap)
ns1.vultr.com
ns2.vultr.com

add DNS in vultr:
A api ip ...

github subdomain:
https://help.github.com/en/github/working-with-github-pages/managing-a-custom-domain-for-your-github-pages-site#configuring-a-subdomain
ip adresses of github pages:
(add in vultr DNS as A with no subdomain)
185.199.108.153
185.199.109.153
185.199.110.153
185.199.111.153

CNAME * skariel.github.io


github pages
-------------------
https://github.com/skariel/td4_front