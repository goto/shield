"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[45],{3905:(e,t,n)=>{n.d(t,{Zo:()=>d,kt:()=>b});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var l=r.createContext({}),p=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},d=function(e){var t=p(e.components);return r.createElement(l.Provider,{value:t},e.children)},u="mdxType",s={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},g=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,l=e.parentName,d=c(e,["components","mdxType","originalType","parentName"]),u=p(n),g=a,b=u["".concat(l,".").concat(g)]||u[g]||s[g]||o;return n?r.createElement(b,i(i({ref:t},d),{},{components:n})):r.createElement(b,i({ref:t},d))}));function b(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=g;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c[u]="string"==typeof e?e:a,i[1]=c;for(var p=2;p<o;p++)i[p]=n[p];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}g.displayName="MDXCreateElement"},5875:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>l,contentTitle:()=>i,default:()=>u,frontMatter:()=>o,metadata:()=>c,toc:()=>p});var r=n(7462),a=(n(7294),n(3905));const o={},i="Creating a group in organization",c={unversionedId:"tour/creating-group",id:"tour/creating-group",title:"Creating a group in organization",description:"In this, we will be using the organization id of the organization we created. Groups in shield belong to one organization.",source:"@site/docs/tour/creating-group.md",sourceDirName:"tour",slug:"/tour/creating-group",permalink:"/shield/tour/creating-group",draft:!1,editUrl:"https://github.com/goto/shield/edit/master/docs/docs/tour/creating-group.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Creating a project in organization",permalink:"/shield/tour/creating-project"},next:{title:"Adding to a group",permalink:"/shield/tour/add-to-group"}},l={},p=[{value:"Relations Table",id:"relations-table",level:3}],d={toc:p};function u(e){let{components:t,...n}=e;return(0,a.kt)("wrapper",(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"creating-a-group-in-organization"},"Creating a group in organization"),(0,a.kt)("p",null,"In this, we will be using the organization id of the organization we created. Groups in shield belong to one organization."),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-sh"},'curl --location --request POST \'http://localhost:8000/admin/v1beta1/groups\'\n--header \'Content-Type: application/json\'\n--data-raw \'{\n    "name": "Data Streaming",\n    "slug": "data-streaming",\n    "metadata": {\n        "description": "group for users in data streaming domain"\n    },\n    "orgId": "4eb3c3b4-962b-4b45-b55b-4c07d3810ca8"\n}\'\n')),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'200\n{\n    "group": {\n        "id": "86e2f95d-92c7-4c59-8fed-b7686cccbf4f",\n        "name": "Data Streaming",\n        "slug": "data-streaming",\n        "orgId": "4eb3c3b4-962b-4b45-b55b-4c07d3810ca8",\n        "metadata": {\n            "description": "group for users in data streaming domain"\n        },\n        "createdAt": "2022-12-07T17:03:59.456847Z",\n        "updatedAt": "2022-12-07T17:03:59.456847Z"\n    }\n}\n')),(0,a.kt)("h3",{id:"relations-table"},"Relations Table"),(0,a.kt)("p",null,"It got an entry for the role ",(0,a.kt)("inlineCode",{parentName:"p"},"group:organization")," for the organization ",(0,a.kt)("inlineCode",{parentName:"p"},"4eb3c3b4-962b-4b45-b55b-4c07d3810ca8"),"."),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-sh"},"                  id                  | subject_namespace_id |              subject_id              | object_namespace_id |              object_id               |        role_id         |          created_at           |          updated_at           | deleted_at \n--------------------------------------+----------------------+--------------------------------------+---------------------+--------------------------------------+------------------------+-------------------------------+-------------------------------+------------\n 460c44a6-f074-4abe-8f8e-949e7a3f5ec2 | user                 | 2fd7f306-61db-4198-9623-6f5f1809df11 | organization        | 4eb3c3b4-962b-4b45-b55b-4c07d3810ca8 | organization:owner     | 2022-12-07 14:10:42.881572+00 | 2022-12-07 14:10:42.881572+00 | \n 10797ec9-6744-4064-8408-c0919e71fbca | organization         | 4eb3c3b4-962b-4b45-b55b-4c07d3810ca8 | project             | 1b89026b-6713-4327-9d7e-ed03345da288 | project:organization   | 2022-12-07 14:31:46.517828+00 | 2022-12-07 14:31:46.517828+00 | \n 29b82d6e-b6fd-4009-9727-1e619c802e23 | organization         | 4eb3c3b4-962b-4b45-b55b-4c07d3810ca8 | group               | 86e2f95d-92c7-4c59-8fed-b7686cccbf4f | group:organization     | 2022-12-07 17:03:59.537254+00 | 2022-12-07 17:03:59.537254+00 |\n(3 rows)\n")))}u.isMDXComponent=!0}}]);