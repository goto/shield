"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[786],{3905:(e,n,t)=>{t.d(n,{Zo:()=>d,kt:()=>g});var a=t(7294);function r(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function o(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);n&&(a=a.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,a)}return t}function i(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?o(Object(t),!0).forEach((function(n){r(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,a,r=function(e,n){if(null==e)return{};var t,a,r={},o=Object.keys(e);for(a=0;a<o.length;a++)t=o[a],n.indexOf(t)>=0||(r[t]=e[t]);return r}(e,n);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)t=o[a],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(r[t]=e[t])}return r}var c=a.createContext({}),s=function(e){var n=a.useContext(c),t=n;return e&&(t="function"==typeof e?e(n):i(i({},n),e)),t},d=function(e){var n=s(e.components);return a.createElement(c.Provider,{value:n},e.children)},p="mdxType",u={inlineCode:"code",wrapper:function(e){var n=e.children;return a.createElement(a.Fragment,{},n)}},m=a.forwardRef((function(e,n){var t=e.components,r=e.mdxType,o=e.originalType,c=e.parentName,d=l(e,["components","mdxType","originalType","parentName"]),p=s(t),m=r,g=p["".concat(c,".").concat(m)]||p[m]||u[m]||o;return t?a.createElement(g,i(i({ref:n},d),{},{components:t})):a.createElement(g,i({ref:n},d))}));function g(e,n){var t=arguments,r=n&&n.mdxType;if("string"==typeof e||r){var o=t.length,i=new Array(o);i[0]=m;var l={};for(var c in n)hasOwnProperty.call(n,c)&&(l[c]=n[c]);l.originalType=e,l[p]="string"==typeof e?e:r,i[1]=l;for(var s=2;s<o;s++)i[s]=t[s];return a.createElement.apply(null,i)}return a.createElement.apply(null,t)}m.displayName="MDXCreateElement"},1541:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>c,contentTitle:()=>i,default:()=>p,frontMatter:()=>o,metadata:()=>l,toc:()=>s});var a=t(7462),r=(t(7294),t(3905));const o={},i="Creating an organization",l={unversionedId:"tour/creating-organization",id:"tour/creating-organization",title:"Creating an organization",description:"Before creating a new organization, let's create an organization admin user.",source:"@site/docs/tour/creating-organization.md",sourceDirName:"tour",slug:"/tour/creating-organization",permalink:"/shield/tour/creating-organization",draft:!1,editUrl:"https://github.com/goto/shield/edit/master/docs/docs/tour/creating-organization.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"What is in Shield?",permalink:"/shield/tour/what-is-in-shield"},next:{title:"Creating a project in organization",permalink:"/shield/tour/creating-project"}},c={},s=[{value:"User creation in Shield",id:"user-creation-in-shield",level:2},{value:"Organization creation in Shield",id:"organization-creation-in-shield",level:2}],d={toc:s};function p(e){let{components:n,...t}=e;return(0,r.kt)("wrapper",(0,a.Z)({},d,t,{components:n,mdxType:"MDXLayout"}),(0,r.kt)("h1",{id:"creating-an-organization"},"Creating an organization"),(0,r.kt)("p",null,"Before creating a new organization, let's create an organization admin user."),(0,r.kt)("h2",{id:"user-creation-in-shield"},"User creation in Shield"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},'curl --location --request POST \'http://localhost:8000/admin/v1beta1/users\'\n--header \'Content-Type: application/json\'\n--header \'X-Shield-Email: admin@gotocompany.com\'\n--data-raw \'{\n    "name": "Shield Org Admin",\n    "email": "admin@gotocompany.com",\n    "metadata": {\n        "role": "organization admin"\n    }\n}\'\n')),(0,r.kt)("p",null,"Note that this will return an error response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'500\n{\n    "code": 13,\n    "message": "internal server error",\n    "details": []\n}\n')),(0,r.kt)("p",null,"This is because metadata key ",(0,r.kt)("inlineCode",{parentName:"p"},"role")," is not defined in ",(0,r.kt)("inlineCode",{parentName:"p"},"metadata_keys")," table. So, let's first create it."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},'curl --location --request POST \'http://localhost:8000/admin/v1beta1/metadatakey\'\n--header \'Content-Type: application/json\'\n--data-raw \'{\n    "key": "role",\n    "description": "role of user in organization"\n}\'\n')),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'200\n{\n    "metadatakey": {\n        "key": "role",\n        "description": "role of user in organization"\n    }\n}\n')),(0,r.kt)("p",null,"Now, we can retry the above user creation request and it should be successful."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'200\n{\n    "user": {\n        "id": "2fd7f306-61db-4198-9623-6f5f1809df11",\n        "name": "Shield Org Admin",\n        "slug": "",\n        "email": "admin@gotocompany.com",\n        "metadata": {\n            "role": "organization admin"\n        },\n        "createdAt": "2022-12-07T13:35:19.005545Z",\n        "updatedAt": "2022-12-07T13:35:19.005545Z"\n    }\n}\n')),(0,r.kt)("p",null,"From now onwards, we can use the above user to perform all the admin operations. Let's begin with organization creation."),(0,r.kt)("h2",{id:"organization-creation-in-shield"},"Organization creation in Shield"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},'curl --location --request POST \'http://localhost:8000/admin/v1beta1/organizations\'\n--header \'Content-Type: application/json\'\n--header \'X-Shield-Email: admin@gotocompany.com\'\n--data-raw \'{\n    "name": "gotocompany",\n    "slug": "gotocompany",\n    "metadata": {\n        "description": "Open DataOps Foundation"\n    }\n}\'\n')),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'200\n{\n    "organization": {\n        "id": "4eb3c3b4-962b-4b45-b55b-4c07d3810ca8",\n        "name": "gotocompany",\n        "slug": "gotocompany",\n        "metadata": {\n            "description": "Open DataOps Foundation"\n        },\n        "createdAt": "2022-12-07T14:10:42.755848Z",\n        "updatedAt": "2022-12-07T14:10:42.755848Z"\n    }\n}\n')),(0,r.kt)("p",null,"Now, let's have a look at relations table where an ",(0,r.kt)("inlineCode",{parentName:"p"},"organization:owner")," relationship is created."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},"                  id                  | subject_namespace_id |              subject_id              | object_namespace_id |              object_id               |      role_id       |          created_at           |          updated_at           | deleted_at \n--------------------------------------+----------------------+--------------------------------------+---------------------+--------------------------------------+--------------------+-------------------------------+-------------------------------+------------\n 460c44a6-f074-4abe-8f8e-949e7a3f5ec2 | user                 | 2fd7f306-61db-4198-9623-6f5f1809df11 | organization        | 4eb3c3b4-962b-4b45-b55b-4c07d3810ca8 | organization:owner | 2022-12-07 14:10:42.881572+00 | 2022-12-07 14:10:42.881572+00 | \n(1 row)\n")))}p.isMDXComponent=!0}}]);