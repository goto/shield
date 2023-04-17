"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[679],{5162:(e,a,t)=>{t.d(a,{Z:()=>r});var n=t(7294),l=t(6010);const s="tabItem_Ymn6";function r(e){let{children:a,hidden:t,className:r}=e;return n.createElement("div",{role:"tabpanel",className:(0,l.Z)(s,r),hidden:t},a)}},5488:(e,a,t)=>{t.d(a,{Z:()=>m});var n=t(7462),l=t(7294),s=t(6010),r=t(2389),o=t(7392),i=t(7094),d=t(2466);const u="tabList__CuJ",c="tabItem_LNqP";function p(e){const{lazy:a,block:t,defaultValue:r,values:p,groupId:m,className:b}=e,h=l.Children.map(e.children,(e=>{if((0,l.isValidElement)(e)&&"value"in e.props)return e;throw new Error(`Docusaurus error: Bad <Tabs> child <${"string"==typeof e.type?e.type:e.type.name}>: all children of the <Tabs> component should be <TabItem>, and every <TabItem> should have a unique "value" prop.`)})),g=p??h.map((e=>{let{props:{value:a,label:t,attributes:n}}=e;return{value:a,label:t,attributes:n}})),k=(0,o.l)(g,((e,a)=>e.value===a.value));if(k.length>0)throw new Error(`Docusaurus error: Duplicate values "${k.map((e=>e.value)).join(", ")}" found in <Tabs>. Every value needs to be unique.`);const T=null===r?r:r??h.find((e=>e.props.default))?.props.value??h[0].props.value;if(null!==T&&!g.some((e=>e.value===T)))throw new Error(`Docusaurus error: The <Tabs> has a defaultValue "${T}" but none of its children has the corresponding value. Available values are: ${g.map((e=>e.value)).join(", ")}. If you intend to show no default tab, use defaultValue={null} instead.`);const{tabGroupChoices:f,setTabGroupChoices:v}=(0,i.U)(),[y,I]=(0,l.useState)(T),C=[],{blockElementScrollPositionUntilNextRender:Z}=(0,d.o5)();if(null!=m){const e=f[m];null!=e&&e!==y&&g.some((a=>a.value===e))&&I(e)}const x=e=>{const a=e.currentTarget,t=C.indexOf(a),n=g[t].value;n!==y&&(Z(a),I(n),null!=m&&v(m,String(n)))},N=e=>{let a=null;switch(e.key){case"Enter":x(e);break;case"ArrowRight":{const t=C.indexOf(e.currentTarget)+1;a=C[t]??C[0];break}case"ArrowLeft":{const t=C.indexOf(e.currentTarget)-1;a=C[t]??C[C.length-1];break}}a?.focus()};return l.createElement("div",{className:(0,s.Z)("tabs-container",u)},l.createElement("ul",{role:"tablist","aria-orientation":"horizontal",className:(0,s.Z)("tabs",{"tabs--block":t},b)},g.map((e=>{let{value:a,label:t,attributes:r}=e;return l.createElement("li",(0,n.Z)({role:"tab",tabIndex:y===a?0:-1,"aria-selected":y===a,key:a,ref:e=>C.push(e),onKeyDown:N,onClick:x},r,{className:(0,s.Z)("tabs__item",c,r?.className,{"tabs__item--active":y===a})}),t??a)}))),a?(0,l.cloneElement)(h.filter((e=>e.props.value===y))[0],{className:"margin-top--md"}):l.createElement("div",{className:"margin-top--md"},h.map(((e,a)=>(0,l.cloneElement)(e,{key:a,hidden:e.props.value!==y})))))}function m(e){const a=(0,r.Z)();return l.createElement(p,(0,n.Z)({key:String(a)},e))}},791:(e,a,t)=>{t.r(a),t.d(a,{assets:()=>c,contentTitle:()=>d,default:()=>b,frontMatter:()=>i,metadata:()=>u,toc:()=>p});var n=t(7462),l=(t(7294),t(3905)),s=t(5488),r=t(5162),o=t(814);const i={},d="Managing Users",u={unversionedId:"guides/managing-user",id:"guides/managing-user",title:"Managing Users",description:"A project in Shield looks like",source:"@site/docs/guides/managing-user.md",sourceDirName:"guides",slug:"/guides/managing-user",permalink:"/shield/guides/managing-user",draft:!1,editUrl:"https://github.com/goto/shield/edit/master/docs/docs/guides/managing-user.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Managing Relations",permalink:"/shield/guides/managing-relation"},next:{title:"Adding Metadata Keys",permalink:"/shield/guides/adding-metadata-key"}},c={},p=[{value:"API Interface",id:"api-interface",level:2},{value:"Create users",id:"create-users",level:3},{value:"List users",id:"list-users",level:3},{value:"Get Users",id:"get-users",level:3},{value:"Update Projects",id:"update-projects",level:3}],m={toc:p};function b(e){let{components:a,...t}=e;return(0,l.kt)("wrapper",(0,n.Z)({},m,t,{components:a,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"managing-users"},"Managing Users"),(0,l.kt)("p",null,"A project in Shield looks like"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "users": [\n        {\n            "id": "598688c6-8c6d-487f-b324-ef3f4af120bb",\n            "name": "John Doe",\n            "slug": "",\n            "email": "john.doe@gotocompany.com",\n            "metadata": {\n                "role": "\\"user-1\\""\n            },\n            "createdAt": "2022-12-09T10:45:19.134019Z",\n            "updatedAt": "2022-12-09T10:45:19.134019Z"\n        }\n    ]\n}\n')),(0,l.kt)("p",null,"One thing to note here is that Shield only allow to have metadata key from a specific set of keys. This constraint is only for users. We can add metadata key using this ",(0,l.kt)("a",{parentName:"p",href:"./adding-metadata-key"},"metadata key API")),(0,l.kt)("h2",{id:"api-interface"},"API Interface"),(0,l.kt)("h3",{id:"create-users"},"Create users"),(0,l.kt)(s.Z,{groupId:"api",mdxType:"Tabs"},(0,l.kt)(r.Z,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.kt)(o.Z,{className:"language-bash",mdxType:"CodeBlock"},'$ curl --location --request POST \'http://localhost:8000/admin/v1beta1/users\'\n--header \'Content-Type: application/json\'\n--header \'Accept: application/json\'\n--header \'X-Shield-Email: admin@gotocompany.com\'\n--data-raw \'{\n  "name": "Jonny Doe",\n  "email": "jonny.doe@gotocompany.com",\n  "metadata": {\n      "role": "user-3"\n  }\n}\'')),(0,l.kt)(r.Z,{value:"CLI",label:"CLI",default:!0,mdxType:"TabItem"},(0,l.kt)(o.Z,{mdxType:"CodeBlock"},(0,l.kt)("p",null,(0,l.kt)("inlineCode",{parentName:"p"},"$ shield user create --file=user.yaml"))))),(0,l.kt)("h3",{id:"list-users"},"List users"),(0,l.kt)(s.Z,{groupId:"api",mdxType:"Tabs"},(0,l.kt)(r.Z,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.kt)(o.Z,{className:"language-bash",mdxType:"CodeBlock"},"curl --location --request GET 'http://localhost:8000/admin/v1beta1/users'\n--header 'Accept: application/json'")),(0,l.kt)(r.Z,{value:"CLI",label:"CLI",default:!0,mdxType:"TabItem"},(0,l.kt)(o.Z,{mdxType:"CodeBlock"},(0,l.kt)("p",null,(0,l.kt)("inlineCode",{parentName:"p"},"$ shield user list"))))),(0,l.kt)("h3",{id:"get-users"},"Get Users"),(0,l.kt)(s.Z,{groupId:"api",mdxType:"Tabs"},(0,l.kt)(r.Z,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.kt)(o.Z,{className:"language-bash",mdxType:"CodeBlock"},"$ curl --location --request GET 'http://localhost:8000/admin/v1beta1/users/e9fba4af-ab23-4631-abba-597b1c8e6608'\n--header 'Accept: application/json''")),(0,l.kt)(r.Z,{value:"CLI",label:"CLI",default:!0,mdxType:"TabItem"},(0,l.kt)(o.Z,{mdxType:"CodeBlock"},(0,l.kt)("p",null,(0,l.kt)("inlineCode",{parentName:"p"},"$ shield user view e9fba4af-ab23-4631-abba-597b1c8e6608 --metadata"))))),(0,l.kt)("h3",{id:"update-projects"},"Update Projects"),(0,l.kt)(s.Z,{groupId:"api",mdxType:"Tabs"},(0,l.kt)(r.Z,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.kt)(o.Z,{className:"language-bash",mdxType:"CodeBlock"},'$ curl --location --request PUT \'http://localhost:8000/admin/v1beta1/users/e9fba4af-ab23-4631-abba-597b1c8e6608\'\n--header \'Content-Type: application/json\'\n--header \'Accept: application/json\'\n--data-raw \'{\n  "name": "Jonny Doe",\n  "email": "john.doe001@gotocompany.com",\n  "metadata": {\n      "role" :   "user-3"\n  }\n}\'')),(0,l.kt)(r.Z,{value:"CLI",label:"CLI",default:!0,mdxType:"TabItem"},(0,l.kt)(o.Z,{mdxType:"CodeBlock"},(0,l.kt)("p",null,(0,l.kt)("inlineCode",{parentName:"p"},"$ shield user edit e9fba4af-ab23-4631-abba-597b1c8e6608 --file=user.yaml"))))))}b.isMDXComponent=!0}}]);