# KamaChat前端
## 项目介绍
kama-chat-server前端，基于vue3+element-plus开发，用于kama-chat-server项目的web端展示。

## 项目结构
```
kama-chat-server前端
├── public
│   ├── favicon.ico
│   └── index.html
├── src
│   ├── assets
│   ├── components
│   ├── router
│   ├── store
│   ├── styles
│   ├── utils
│   ├── views
│   ├── App.vue
│   ├── main.ts
```


## 项目运行
```shell
vue create kama-chat-server --registry=https://registry.npmmirror.com
```
vue-cli创建项目，选择默认配置，项目名称为kama-chat-server，使用npm镜像源。


```
Vue CLI v5.0.8
? Please pick a preset: Manually select features
? Check the features needed for your project: Babel, Router, Vuex, Linter 
? Choose a version of Vue.js that you want to start the project with 3.x
? Use history mode for router? (Requires proper server setup for index 
fallback in production) Yes
? Pick a linter / formatter config: Basic
? Pick additional lint features: Lint on save
? Where do you prefer placing config for Babel, ESLint, etc.? In dedicated 
config files
? Save this as a preset for future projects? No
```
以上为创建项目时的配置，跟着选就行了。

```shell
yarn add element-plus -g
yarn add axios -g
yarn add @element-plus/icons-vue -g
```
配置element-plus，axios，element-plus的图标等。

```shell
yarn serve
```
启动项目。


## 参考文献
https://www.iconfont.cn/ 阿里图标库
