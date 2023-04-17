"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[266],{3905:(e,n,a)=>{a.d(n,{Zo:()=>g,kt:()=>m});var t=a(7294);function i(e,n,a){return n in e?Object.defineProperty(e,n,{value:a,enumerable:!0,configurable:!0,writable:!0}):e[n]=a,e}function r(e,n){var a=Object.keys(e);if(Object.getOwnPropertySymbols){var t=Object.getOwnPropertySymbols(e);n&&(t=t.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),a.push.apply(a,t)}return a}function o(e){for(var n=1;n<arguments.length;n++){var a=null!=arguments[n]?arguments[n]:{};n%2?r(Object(a),!0).forEach((function(n){i(e,n,a[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(a)):r(Object(a)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(a,n))}))}return e}function s(e,n){if(null==e)return{};var a,t,i=function(e,n){if(null==e)return{};var a,t,i={},r=Object.keys(e);for(t=0;t<r.length;t++)a=r[t],n.indexOf(a)>=0||(i[a]=e[a]);return i}(e,n);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);for(t=0;t<r.length;t++)a=r[t],n.indexOf(a)>=0||Object.prototype.propertyIsEnumerable.call(e,a)&&(i[a]=e[a])}return i}var l=t.createContext({}),c=function(e){var n=t.useContext(l),a=n;return e&&(a="function"==typeof e?e(n):o(o({},n),e)),a},g=function(e){var n=c(e.components);return t.createElement(l.Provider,{value:n},e.children)},p="mdxType",u={inlineCode:"code",wrapper:function(e){var n=e.children;return t.createElement(t.Fragment,{},n)}},d=t.forwardRef((function(e,n){var a=e.components,i=e.mdxType,r=e.originalType,l=e.parentName,g=s(e,["components","mdxType","originalType","parentName"]),p=c(a),d=i,m=p["".concat(l,".").concat(d)]||p[d]||u[d]||r;return a?t.createElement(m,o(o({ref:n},g),{},{components:a})):t.createElement(m,o({ref:n},g))}));function m(e,n){var a=arguments,i=n&&n.mdxType;if("string"==typeof e||i){var r=a.length,o=new Array(r);o[0]=d;var s={};for(var l in n)hasOwnProperty.call(n,l)&&(s[l]=n[l]);s.originalType=e,s[p]="string"==typeof e?e:i,o[1]=s;for(var c=2;c<r;c++)o[c]=a[c];return t.createElement.apply(null,o)}return t.createElement.apply(null,a)}d.displayName="MDXCreateElement"},4967:(e,n,a)=>{a.r(n),a.d(n,{assets:()=>l,contentTitle:()=>o,default:()=>p,frontMatter:()=>r,metadata:()=>s,toc:()=>c});var t=a(7462),i=(a(7294),a(3905));const r={},o="Overview",s={unversionedId:"guides/overview",id:"guides/overview",title:"Overview",description:"The following topics will describe how to use Shield. It respects multi-tenancy using namespace, which can either be a system namespace or a resource namespace.",source:"@site/docs/guides/overview.md",sourceDirName:"guides",slug:"/guides/overview",permalink:"/shield/guides/overview",draft:!1,editUrl:"https://github.com/goto/shield/edit/master/docs/docs/guides/overview.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Shield as a proxy",permalink:"/shield/tour/shield-as-proxy"},next:{title:"Managing Organization",permalink:"/shield/guides/managing-organization"}},l={},c=[{value:"Managing Organizations",id:"managing-organizations",level:2},{value:"Managing Projects",id:"managing-projects",level:2},{value:"Managing Resources",id:"managing-resources",level:2},{value:"Managing Groups",id:"managing-groups",level:2},{value:"Managing Namespaces, Policies, Roles and Actions",id:"managing-namespaces-policies-roles-and-actions",level:2},{value:"Managing Users and their Metadata",id:"managing-users-and-their-metadata",level:2},{value:"Managing Relations",id:"managing-relations",level:2},{value:"Checking Permission",id:"checking-permission",level:2},{value:"Where to go next?",id:"where-to-go-next",level:2}],g={toc:c};function p(e){let{components:n,...a}=e;return(0,i.kt)("wrapper",(0,t.Z)({},g,a,{components:n,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"overview"},"Overview"),(0,i.kt)("p",null,"The following topics will describe how to use Shield. It respects multi-tenancy using namespace, which can either be a system namespace or a resource namespace.\nSystem namespace is composed of a ",(0,i.kt)("inlineCode",{parentName:"p"},"backend")," and a ",(0,i.kt)("inlineCode",{parentName:"p"},"resource type")," which allows onboard multiple instances of a service by changing the backend. While resource namespaces allows to onboard multiple ",(0,i.kt)("inlineCode",{parentName:"p"},"organizations"),", ",(0,i.kt)("inlineCode",{parentName:"p"},"project")," and ",(0,i.kt)("inlineCode",{parentName:"p"},"groups"),"."),(0,i.kt)("h2",{id:"managing-organizations"},"Managing Organizations"),(0,i.kt)("p",null,"Organizations are the top most level object in Sheild's system."),(0,i.kt)("h2",{id:"managing-projects"},"Managing Projects"),(0,i.kt)("p",null,"Project comes under an organization, and they can have multiple resources belonging to them."),(0,i.kt)("h2",{id:"managing-resources"},"Managing Resources"),(0,i.kt)("p",null,"Resources have some basic information about the resources being created on the backend. They can be used for authorization purpose."),(0,i.kt)("h2",{id:"managing-groups"},"Managing Groups"),(0,i.kt)("p",null,"Groups fall under anorganization too, an they are a colection of users with different roles."),(0,i.kt)("h2",{id:"managing-namespaces-policies-roles-and-actions"},"Managing Namespaces, Policies, Roles and Actions"),(0,i.kt)("p",null,"All of these are managed by configuration files and shall not be modified via APIs."),(0,i.kt)("h2",{id:"managing-users-and-their-metadata"},"Managing Users and their Metadata"),(0,i.kt)("p",null,"User represent a real life user distinguished by their emails."),(0,i.kt)("h2",{id:"managing-relations"},"Managing Relations"),(0,i.kt)("p",null,"Relations are a copy of the relationships being managed in SpiceDB."),(0,i.kt)("h2",{id:"checking-permission"},"Checking Permission"),(0,i.kt)("p",null,"Shield provides an API to check if a user has a  certain permissions on a resource."),(0,i.kt)("h2",{id:"where-to-go-next"},"Where to go next?"),(0,i.kt)("p",null,"We recomment you to check all the guides for having a clear understanding of the APIs. For testing these APIs on local, you can import the ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/goto/shield/blob/main/proto/apidocs.swagger.json"},"Swagger"),"."))}p.isMDXComponent=!0}}]);