"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[217],{3905:(e,t,r)=>{r.d(t,{Zo:()=>p,kt:()=>h});var n=r(7294);function l(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){l(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function a(e,t){if(null==e)return{};var r,n,l=function(e,t){if(null==e)return{};var r,n,l={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(l[r]=e[r]);return l}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(l[r]=e[r])}return l}var s=n.createContext({}),c=function(e){var t=n.useContext(s),r=t;return e&&(r="function"==typeof e?e(t):i(i({},t),e)),r},p=function(e){var t=c(e.components);return n.createElement(s.Provider,{value:t},e.children)},d="mdxType",u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},m=n.forwardRef((function(e,t){var r=e.components,l=e.mdxType,o=e.originalType,s=e.parentName,p=a(e,["components","mdxType","originalType","parentName"]),d=c(r),m=l,h=d["".concat(s,".").concat(m)]||d[m]||u[m]||o;return r?n.createElement(h,i(i({ref:t},p),{},{components:r})):n.createElement(h,i({ref:t},p))}));function h(e,t){var r=arguments,l=t&&t.mdxType;if("string"==typeof e||l){var o=r.length,i=new Array(o);i[0]=m;var a={};for(var s in t)hasOwnProperty.call(t,s)&&(a[s]=t[s]);a.originalType=e,a[d]="string"==typeof e?e:l,i[1]=a;for(var c=2;c<o;c++)i[c]=r[c];return n.createElement.apply(null,i)}return n.createElement.apply(null,r)}m.displayName="MDXCreateElement"},9803:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>s,contentTitle:()=>i,default:()=>d,frontMatter:()=>o,metadata:()=>a,toc:()=>c});var n=r(7462),l=(r(7294),r(3905));const o={},i="Installation",a={unversionedId:"installation",id:"installation",title:"Installation",description:"We provide pre-built binaries, Docker Images and Helm Charts",source:"@site/docs/installation.md",sourceDirName:".",slug:"/installation",permalink:"/shield/installation",draft:!1,editUrl:"https://github.com/odpf/shield/edit/master/docs/docs/installation.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Introduction",permalink:"/shield/"},next:{title:"Architecture",permalink:"/shield/concepts/architecture"}},s={},c=[{value:"Binary (Cross-platform)",id:"binary-cross-platform",level:2},{value:"Homebrew",id:"homebrew",level:2},{value:"Docker",id:"docker",level:2},{value:"Prerequisites",id:"prerequisites",level:3},{value:"Compiling from source",id:"compiling-from-source",level:2},{value:"Prerequisites",id:"prerequisites-1",level:3},{value:"Build",id:"build",level:3}],p={toc:c};function d(e){let{components:t,...r}=e;return(0,l.kt)("wrapper",(0,n.Z)({},p,r,{components:t,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"installation"},"Installation"),(0,l.kt)("p",null,"We provide pre-built ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/odpf/shield/releases"},"binaries"),", ",(0,l.kt)("a",{parentName:"p",href:"https://hub.docker.com/r/odpf/shield"},"Docker Images")," and ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/odpf/charts/tree/main/stable/shield"},"Helm Charts")),(0,l.kt)("h2",{id:"binary-cross-platform"},"Binary (Cross-platform)"),(0,l.kt)("p",null,"Download the appropriate version for your platform from ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/odpf/shield/releases"},"releases")," page. Once downloaded, the binary can be run from anywhere.\nYou don\u2019t need to install it into a global location. This works well for shared hosts and other systems where you don\u2019t have a privileged account.\nIdeally, you should install it somewhere in your PATH for easy use. ",(0,l.kt)("inlineCode",{parentName:"p"},"/usr/local/bin")," is the most probable location."),(0,l.kt)("h2",{id:"homebrew"},"Homebrew"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-sh"},"# Install shield (requires homebrew installed)\n$ brew install odpf/taps/shield\n\n# Upgrade shield (requires homebrew installed)\n$ brew upgrade shield\n\n# Check for installed shield version\n$ shield version\n")),(0,l.kt)("h2",{id:"docker"},"Docker"),(0,l.kt)("h3",{id:"prerequisites"},"Prerequisites"),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},"Docker installed")),(0,l.kt)("p",null,"Run Docker Image"),(0,l.kt)("p",null,"Shield provides Docker image as part of the release. Make sure you have Spicedb and postgres running on your local and run the following."),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-sh"},"# Download docker image from docker hub\n$ docker pull odpf/shield\n\n# Run the following docker command with minimal config.\n$ docker run -p 8080:8080 \\\n  -e SHIELD_DB_DRIVER=postgres \\\n  -e SHIELD_DB_URL=postgres://shield:@localhost:5432/shield?sslmode=disable \\\n  -e SHIELD_SPICEDB_HOST=spicedb.localhost:50051 \\\n  -e SHIELD_SPICEDB_PRE_SHARED_KEY=randomkey\n  -v .config:.config\n  odpf/shield serve\n")),(0,l.kt)("h2",{id:"compiling-from-source"},"Compiling from source"),(0,l.kt)("h3",{id:"prerequisites-1"},"Prerequisites"),(0,l.kt)("p",null,"Shield requires the following dependencies:"),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},"Golang (version 1.18 or above)"),(0,l.kt)("li",{parentName:"ul"},"Git")),(0,l.kt)("h3",{id:"build"},"Build"),(0,l.kt)("p",null,"Run the following commands to compile ",(0,l.kt)("inlineCode",{parentName:"p"},"shield")," from source"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-shell"},"git clone git@github.com:odpf/shield.git\ncd shield\nmake build\n")),(0,l.kt)("p",null,"Use the following command to test"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-shell"},"./shield version\n")),(0,l.kt)("p",null,"Shield service can be started with the following command although there are few required ",(0,l.kt)("a",{parentName:"p",href:"/shield/reference/configurations"},"configurations")," for it to start."),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-sh"},"./shield server start\n")))}d.isMDXComponent=!0}}]);