## 设置wsl的IP到windows宿主机的hosts文件中的小工具

自动读取wsl的ip地址，并添加或替换到hosts文件内，如需要开机启动可添加到启动菜单。

### 使用方法：

#### 方法1：将需要配置的域名写入到同目录的 %s 文件中，运行将自动按行读取

#### 方法2：将域名追加到运行命令后

例如： ./wsl-host-set.exe local-website.com buyaoqiao.tech liuggchen.com
