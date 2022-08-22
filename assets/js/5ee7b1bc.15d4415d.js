"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[810],{3905:function(e,r,t){t.d(r,{Zo:function(){return l},kt:function(){return m}});var a=t(7294);function n(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function o(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);r&&(a=a.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,a)}return t}function i(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?o(Object(t),!0).forEach((function(r){n(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function s(e,r){if(null==e)return{};var t,a,n=function(e,r){if(null==e)return{};var t,a,n={},o=Object.keys(e);for(a=0;a<o.length;a++)t=o[a],r.indexOf(t)>=0||(n[t]=e[t]);return n}(e,r);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)t=o[a],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(n[t]=e[t])}return n}var c=a.createContext({}),p=function(e){var r=a.useContext(c),t=r;return e&&(t="function"==typeof e?e(r):i(i({},r),e)),t},l=function(e){var r=p(e.components);return a.createElement(c.Provider,{value:r},e.children)},u={inlineCode:"code",wrapper:function(e){var r=e.children;return a.createElement(a.Fragment,{},r)}},d=a.forwardRef((function(e,r){var t=e.components,n=e.mdxType,o=e.originalType,c=e.parentName,l=s(e,["components","mdxType","originalType","parentName"]),d=p(t),m=n,f=d["".concat(c,".").concat(m)]||d[m]||u[m]||o;return t?a.createElement(f,i(i({ref:r},l),{},{components:t})):a.createElement(f,i({ref:r},l))}));function m(e,r){var t=arguments,n=r&&r.mdxType;if("string"==typeof e||n){var o=t.length,i=new Array(o);i[0]=d;var s={};for(var c in r)hasOwnProperty.call(r,c)&&(s[c]=r[c]);s.originalType=e,s.mdxType="string"==typeof e?e:n,i[1]=s;for(var p=2;p<o;p++)i[p]=t[p];return a.createElement.apply(null,i)}return a.createElement.apply(null,t)}d.displayName="MDXCreateElement"},1480:function(e,r,t){t.r(r),t.d(r,{assets:function(){return l},contentTitle:function(){return c},default:function(){return m},frontMatter:function(){return s},metadata:function(){return p},toc:function(){return u}});var a=t(7462),n=t(3366),o=(t(7294),t(3905)),i=["components"],s={},c="Glossary",p={unversionedId:"concepts/glossary",id:"concepts/glossary",title:"Glossary",description:"User",source:"@site/docs/concepts/glossary.md",sourceDirName:"concepts",slug:"/concepts/glossary",permalink:"/shield/concepts/glossary",draft:!1,editUrl:"https://github.com/odpf/shield/edit/master/docs/docs/concepts/glossary.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Architecture",permalink:"/shield/concepts/architecture"},next:{title:"Configurations",permalink:"/shield/reference/configurations"}},l={},u=[{value:"User",id:"user",level:2},{value:"Group",id:"group",level:2},{value:"Role",id:"role",level:2},{value:"Namespace",id:"namespace",level:2},{value:"Policy",id:"policy",level:2},{value:"Action",id:"action",level:2},{value:"Project",id:"project",level:2},{value:"Orgnaization",id:"orgnaization",level:2},{value:"Resource",id:"resource",level:2},{value:"Spicedb",id:"spicedb",level:2}],d={toc:u};function m(e){var r=e.components,t=(0,n.Z)(e,i);return(0,o.kt)("wrapper",(0,a.Z)({},d,t,{components:r,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"glossary"},"Glossary"),(0,o.kt)("h2",{id:"user"},"User"),(0,o.kt)("p",null,"A User represents a profile of a user. User has to create a profile with their name, email, and metadata."),(0,o.kt)("h2",{id:"group"},"Group"),(0,o.kt)("p",null,"A Group represents a group of ",(0,o.kt)("a",{parentName:"p",href:"#user"},"users"),". There are predefined ",(0,o.kt)("a",{parentName:"p",href:"#role"},"roles")," for a user like Team Admin and Team Members. By default, every user gets a Team member role. Team admin can add and remove members from the Group."),(0,o.kt)("h2",{id:"role"},"Role"),(0,o.kt)("p",null,"Within a ",(0,o.kt)("a",{parentName:"p",href:"#namespace"},"namespace"),", a role will get assigned to ",(0,o.kt)("a",{parentName:"p",href:"#user"},"users"),". A role is needed to define ",(0,o.kt)("a",{parentName:"p",href:"#policy"},"policy"),". There are predefined roles in a namespace. Users can also create custom roles."),(0,o.kt)("h2",{id:"namespace"},"Namespace"),(0,o.kt)("p",null,"Namespace provides scope for the ",(0,o.kt)("a",{parentName:"p",href:"#policy"},"policies"),". There are predefined namespaces like ",(0,o.kt)("a",{parentName:"p",href:"#orgnaization"},"organization"),", ",(0,o.kt)("a",{parentName:"p",href:"#project"},"project"),", ",(0,o.kt)("a",{parentName:"p",href:"#group"},"group"),", etc. Users can also create custom namespaces."),(0,o.kt)("h2",{id:"policy"},"Policy"),(0,o.kt)("p",null,"A Policy defines what ",(0,o.kt)("a",{parentName:"p",href:"#action"},"actions")," a ",(0,o.kt)("a",{parentName:"p",href:"#role"},"role")," can perform in a ",(0,o.kt)("a",{parentName:"p",href:"#namespace"},"namespace"),". All policies in a namespace will be used to generate schema for ",(0,o.kt)("a",{parentName:"p",href:"#spicedb"},"spicedb"),"."),(0,o.kt)("h2",{id:"action"},"Action"),(0,o.kt)("p",null,"Within a ",(0,o.kt)("a",{parentName:"p",href:"#namespace"},"namespace"),", a role can perform certain actions. A action is needed to define ",(0,o.kt)("a",{parentName:"p",href:"#policy"},"policy"),"."),(0,o.kt)("h2",{id:"project"},"Project"),(0,o.kt)("p",null,"A Project is a scope in which a User can create and manage Resources."),(0,o.kt)("h2",{id:"orgnaization"},"Orgnaization"),(0,o.kt)("p",null,"An Orgnaization is just a group of Projects. ",(0,o.kt)("a",{parentName:"p",href:"#group"},"Groups")," also belongs to Orgnaization."),(0,o.kt)("h2",{id:"resource"},"Resource"),(0,o.kt)("p",null,"Resources are just custom namespaces in which user can define policies. All Resources also gets default policies."),(0,o.kt)("h2",{id:"spicedb"},"Spicedb"),(0,o.kt)("p",null,(0,o.kt)("a",{parentName:"p",href:"https://github.com/authzed/spicedb"},"SpiceDB")," is a Zanzibar-inspired open source database system for managing security-critical application permissions."))}m.isMDXComponent=!0}}]);