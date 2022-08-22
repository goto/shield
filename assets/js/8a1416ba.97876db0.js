"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[886],{3905:function(e,t,r){r.d(t,{Zo:function(){return p},kt:function(){return h}});var n=r(7294);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function a(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function c(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},i=Object.keys(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var s=n.createContext({}),l=function(e){var t=n.useContext(s),r=t;return e&&(r="function"==typeof e?e(t):a(a({},t),e)),r},p=function(e){var t=l(e.components);return n.createElement(s.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},d=n.forwardRef((function(e,t){var r=e.components,o=e.mdxType,i=e.originalType,s=e.parentName,p=c(e,["components","mdxType","originalType","parentName"]),d=l(r),h=o,f=d["".concat(s,".").concat(h)]||d[h]||u[h]||i;return r?n.createElement(f,a(a({ref:t},p),{},{components:r})):n.createElement(f,a({ref:t},p))}));function h(e,t){var r=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=r.length,a=new Array(i);a[0]=d;var c={};for(var s in t)hasOwnProperty.call(t,s)&&(c[s]=t[s]);c.originalType=e,c.mdxType="string"==typeof e?e:o,a[1]=c;for(var l=2;l<i;l++)a[l]=r[l];return n.createElement.apply(null,a)}return n.createElement.apply(null,r)}d.displayName="MDXCreateElement"},4730:function(e,t,r){r.r(t),r.d(t,{assets:function(){return p},contentTitle:function(){return s},default:function(){return h},frontMatter:function(){return c},metadata:function(){return l},toc:function(){return u}});var n=r(7462),o=r(3366),i=(r(7294),r(3905)),a=["components"],c={},s="Architecture",l={unversionedId:"concepts/architecture",id:"concepts/architecture",title:"Architecture",description:"Shield exposes both HTTP and gRPC APIs to manage data. It also proxy APIs to other services. Shield talks to SpiceDB instance to check for authorization.",source:"@site/docs/concepts/architecture.md",sourceDirName:"concepts",slug:"/concepts/architecture",permalink:"/shield/concepts/architecture",draft:!1,editUrl:"https://github.com/odpf/shield/edit/master/docs/docs/concepts/architecture.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Manage Resources",permalink:"/shield/guides/manage-resources"},next:{title:"Glossary",permalink:"/shield/concepts/glossary"}},p={},u=[{value:"Technologies",id:"technologies",level:2},{value:"Components",id:"components",level:2},{value:"API and Proxy Server",id:"api-and-proxy-server",level:3},{value:"PostgresDB",id:"postgresdb",level:3},{value:"SpiceDB",id:"spicedb",level:3}],d={toc:u};function h(e){var t=e.components,c=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,n.Z)({},d,c,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"architecture"},"Architecture"),(0,i.kt)("p",null,"Shield exposes both HTTP and gRPC APIs to manage data. It also proxy APIs to other services. Shield talks to SpiceDB instance to check for authorization."),(0,i.kt)("p",null,(0,i.kt)("img",{alt:"Shield Architecture",src:r(6756).Z,width:"1292",height:"649"})),(0,i.kt)("h2",{id:"technologies"},"Technologies"),(0,i.kt)("p",null,"Shield is developed with"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Golang - Programming language"),(0,i.kt)("li",{parentName:"ul"},"Docker - container engine to start postgres and cortex to aid development"),(0,i.kt)("li",{parentName:"ul"},"Postgres - a relational database"),(0,i.kt)("li",{parentName:"ul"},"SpiceDB - SpiceDB is an open source database system for managing security-critical application permissions.")),(0,i.kt)("h2",{id:"components"},"Components"),(0,i.kt)("h3",{id:"api-and-proxy-server"},"API and Proxy Server"),(0,i.kt)("p",null,"Shield server exposes both HTTP and gRPC APIs (via GRPC gateway) to manage users, groups, policies, etc. It also runs a proxy server on different port."),(0,i.kt)("h3",{id:"postgresdb"},"PostgresDB"),(0,i.kt)("p",null,"There are 2 PostgresDB instances. One instance is required for Shield to store all the business logic like user detail, team detail, User's role in the team, etc."),(0,i.kt)("p",null,"Another DB instance is for SpiceDB to store all the data needed for authorization."),(0,i.kt)("h3",{id:"spicedb"},"SpiceDB"),(0,i.kt)("p",null,"Shield push all the policies and relationships data to SpiceDB. All this data is needed to make the authorization decision. Shield connects to SpiceDB instance via gRPC"))}h.isMDXComponent=!0},6756:function(e,t,r){t.Z=r.p+"assets/images/architecture-f1d636d19cae9540910a20266e2c8179.svg"}}]);