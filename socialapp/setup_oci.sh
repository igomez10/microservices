sudp iptables -I INPUT -j ACCEPT
sudo su -
iptables-legacy-save > /etc/iptables/rules.v4
exit

# install docker
sudo apt update
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository -y "deb [arch=arm64] https://download.docker.com/linux/ubuntu focal stable"
apt-cache policy docker-ce
sudo apt install -y docker-ce
sudo systemctl status docker
sudo usermod -aG docker ${USER}
groups
sudo usermod -aG docker ubuntu
sudo chmod 666 /var/run/docker.sock

# install go
wget https://go.dev/dl/go1.19.2.linux-arm64.tar.gz
rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.19.2.linux-arm64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# install terraform
sudo apt install -y unzip
wget https://releases.hashicorp.com/terraform/1.3.3/terraform_1.3.3_linux_arm64.zip
unzip terraform_1.3.3_linux_arm64.zip
sudo mv terraform /usr/local/bin/

# install oh my zsh
sudo apt install -y zsh
sh -c "$(wget https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh -O -)"
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k
# Set ZSH_THEME="powerlevel10k/powerlevel10k" in ~/.zshrc.
echo "ZSH_THEME=\"powerlevel10k/powerlevel10k\"" >> ~/.zshrc
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ~/powerlevel10k
echo 'source ~/powerlevel10k/powerlevel10k.zsh-theme' >>~/.zshrc
#  Set shell to zsh
sudo chsh -s $(which zsh)
## syntax highlighting
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git
echo "source ${(q-)PWD}/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh" >> ${ZDOTDIR:-$HOME}/.zshrc
## completions
git clone https://github.com/zsh-users/zsh-completions ${ZSH_CUSTOM:-${ZSH:-~/.oh-my-zsh}/custom}/plugins/zsh-completions
## auto suggestions
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions

## silver searcher
sudo apt install -y silversearcher-ag
## bat
sudo apt install -y bat
## fzf
sudo apt-get install -y fzf
## tldr
sudo apt install -y tldr
## ncdu
sudo apt install -y ncdu
## docker compose
sudo apt install -y docker-compose
