"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[569],{9365:(e,a,t)=>{t.d(a,{A:()=>r});var n=t(6540),l=t(53);const s="tabItem_Ymn6";function r(e){let{children:a,hidden:t,className:r}=e;return n.createElement("div",{role:"tabpanel",className:(0,l.A)(s,r),hidden:t},a)}},4865:(e,a,t)=>{t.d(a,{A:()=>m});var n=t(8168),l=t(6540),s=t(53),r=t(2303),o=t(1682),i=t(6976),d=t(3104);const u="tabList__CuJ",c="tabItem_LNqP";function p(e){const{lazy:a,block:t,defaultValue:r,values:p,groupId:m,className:g}=e,y=l.Children.map(e.children,(e=>{if((0,l.isValidElement)(e)&&"value"in e.props)return e;throw new Error(`Docusaurus error: Bad <Tabs> child <${"string"==typeof e.type?e.type:e.type.name}>: all children of the <Tabs> component should be <TabItem>, and every <TabItem> should have a unique "value" prop.`)})),b=p??y.map((e=>{let{props:{value:a,label:t,attributes:n}}=e;return{value:a,label:t,attributes:n}})),h=(0,o.X)(b,((e,a)=>e.value===a.value));if(h.length>0)throw new Error(`Docusaurus error: Duplicate values "${h.map((e=>e.value)).join(", ")}" found in <Tabs>. Every value needs to be unique.`);const T=null===r?r:r??y.find((e=>e.props.default))?.props.value??y[0].props.value;if(null!==T&&!b.some((e=>e.value===T)))throw new Error(`Docusaurus error: The <Tabs> has a defaultValue "${T}" but none of its children has the corresponding value. Available values are: ${b.map((e=>e.value)).join(", ")}. If you intend to show no default tab, use defaultValue={null} instead.`);const{tabGroupChoices:f,setTabGroupChoices:v}=(0,i.x)(),[A,k]=(0,l.useState)(T),I=[],{blockElementScrollPositionUntilNextRender:C}=(0,d.a_)();if(null!=m){const e=f[m];null!=e&&e!==A&&b.some((a=>a.value===e))&&k(e)}const x=e=>{const a=e.currentTarget,t=I.indexOf(a),n=b[t].value;n!==A&&(C(a),k(n),null!=m&&v(m,String(n)))},N=e=>{let a=null;switch(e.key){case"Enter":x(e);break;case"ArrowRight":{const t=I.indexOf(e.currentTarget)+1;a=I[t]??I[0];break}case"ArrowLeft":{const t=I.indexOf(e.currentTarget)-1;a=I[t]??I[I.length-1];break}}a?.focus()};return l.createElement("div",{className:(0,s.A)("tabs-container",u)},l.createElement("ul",{role:"tablist","aria-orientation":"horizontal",className:(0,s.A)("tabs",{"tabs--block":t},g)},b.map((e=>{let{value:a,label:t,attributes:r}=e;return l.createElement("li",(0,n.A)({role:"tab",tabIndex:A===a?0:-1,"aria-selected":A===a,key:a,ref:e=>I.push(e),onKeyDown:N,onClick:x},r,{className:(0,s.A)("tabs__item",c,r?.className,{"tabs__item--active":A===a})}),t??a)}))),a?(0,l.cloneElement)(y.filter((e=>e.props.value===A))[0],{className:"margin-top--md"}):l.createElement("div",{className:"margin-top--md"},y.map(((e,a)=>(0,l.cloneElement)(e,{key:a,hidden:e.props.value!==A})))))}function m(e){const a=(0,r.A)();return l.createElement(p,(0,n.A)({key:String(a)},e))}},6926:(e,a,t)=>{t.r(a),t.d(a,{assets:()=>c,contentTitle:()=>d,default:()=>g,frontMatter:()=>i,metadata:()=>u,toc:()=>p});var n=t(8168),l=(t(6540),t(5680)),s=t(4865),r=t(9365),o=t(7964);const i={},d="Managing Users",u={unversionedId:"guides/managing-user",id:"guides/managing-user",title:"Managing Users",description:"A project in Shield looks like",source:"@site/docs/guides/managing-user.md",sourceDirName:"guides",slug:"/guides/managing-user",permalink:"/shield/guides/managing-user",draft:!1,editUrl:"https://github.com/goto/shield/edit/master/docs/docs/guides/managing-user.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Managing Relations",permalink:"/shield/guides/managing-relation"},next:{title:"Adding Metadata Keys",permalink:"/shield/guides/adding-metadata-key"}},c={},p=[{value:"API Interface",id:"api-interface",level:2},{value:"Create users",id:"create-users",level:3},{value:"List users",id:"list-users",level:3},{value:"Get Users",id:"get-users",level:3},{value:"Update Projects",id:"update-projects",level:3}],m={toc:p};function g(e){let{components:a,...t}=e;return(0,l.yg)("wrapper",(0,n.A)({},m,t,{components:a,mdxType:"MDXLayout"}),(0,l.yg)("h1",{id:"managing-users"},"Managing Users"),(0,l.yg)("p",null,"A project in Shield looks like"),(0,l.yg)("pre",null,(0,l.yg)("code",{parentName:"pre",className:"language-json"},'{\n    "users": [\n        {\n            "id": "598688c6-8c6d-487f-b324-ef3f4af120bb",\n            "name": "John Doe",\n            "slug": "",\n            "email": "john.doe@gotocompany.com",\n            "metadata": {\n                "role": "\\"user-1\\""\n            },\n            "createdAt": "2022-12-09T10:45:19.134019Z",\n            "updatedAt": "2022-12-09T10:45:19.134019Z"\n        }\n    ]\n}\n')),(0,l.yg)("p",null,"One thing to note here is that Shield only allow to have metadata key from a specific set of keys. This constraint is only for users. We can add metadata key using this ",(0,l.yg)("a",{parentName:"p",href:"./adding-metadata-key"},"metadata key API")),(0,l.yg)("h2",{id:"api-interface"},"API Interface"),(0,l.yg)("h3",{id:"create-users"},"Create users"),(0,l.yg)(s.A,{groupId:"api",mdxType:"Tabs"},(0,l.yg)(r.A,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.yg)(o.A,{className:"language-bash",mdxType:"CodeBlock"},'$ curl --location --request POST \'http://localhost:8000/admin/v1beta1/users\'\n--header \'Content-Type: application/json\'\n--header \'Accept: application/json\'\n--header \'X-Shield-Email: admin@gotocompany.com\'\n--data-raw \'{\n  "name": "Jonny Doe",\n  "email": "jonny.doe@gotocompany.com",\n  "metadata": {\n      "role": "user-3"\n  }\n}\'')),(0,l.yg)(r.A,{value:"CLI",label:"CLI",default:!0,mdxType:"TabItem"},(0,l.yg)(o.A,{mdxType:"CodeBlock"},(0,l.yg)("p",null,(0,l.yg)("inlineCode",{parentName:"p"},"$ shield user create --file=user.yaml"))))),(0,l.yg)("h3",{id:"list-users"},"List users"),(0,l.yg)(s.A,{groupId:"api",mdxType:"Tabs"},(0,l.yg)(r.A,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.yg)(o.A,{className:"language-bash",mdxType:"CodeBlock"},"curl --location --request GET 'http://localhost:8000/admin/v1beta1/users'\n--header 'Accept: application/json'")),(0,l.yg)(r.A,{value:"CLI",label:"CLI",default:!0,mdxType:"TabItem"},(0,l.yg)(o.A,{mdxType:"CodeBlock"},(0,l.yg)("p",null,(0,l.yg)("inlineCode",{parentName:"p"},"$ shield user list"))))),(0,l.yg)("h3",{id:"get-users"},"Get Users"),(0,l.yg)(s.A,{groupId:"api",mdxType:"Tabs"},(0,l.yg)(r.A,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.yg)(o.A,{className:"language-bash",mdxType:"CodeBlock"},"$ curl --location --request GET 'http://localhost:8000/admin/v1beta1/users/e9fba4af-ab23-4631-abba-597b1c8e6608'\n--header 'Accept: application/json''")),(0,l.yg)(r.A,{value:"CLI",label:"CLI",default:!0,mdxType:"TabItem"},(0,l.yg)(o.A,{mdxType:"CodeBlock"},(0,l.yg)("p",null,(0,l.yg)("inlineCode",{parentName:"p"},"$ shield user view e9fba4af-ab23-4631-abba-597b1c8e6608 --metadata"))))),(0,l.yg)("h3",{id:"update-projects"},"Update Projects"),(0,l.yg)(s.A,{groupId:"api",mdxType:"Tabs"},(0,l.yg)(r.A,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.yg)(o.A,{className:"language-bash",mdxType:"CodeBlock"},'$ curl --location --request PUT \'http://localhost:8000/admin/v1beta1/users/e9fba4af-ab23-4631-abba-597b1c8e6608\'\n--header \'Content-Type: application/json\'\n--header \'Accept: application/json\'\n--data-raw \'{\n  "name": "Jonny Doe",\n  "email": "john.doe001@gotocompany.com",\n  "metadata": {\n      "role" :   "user-3"\n  }\n}\'')),(0,l.yg)(r.A,{value:"CLI",label:"CLI",default:!0,mdxType:"TabItem"},(0,l.yg)(o.A,{mdxType:"CodeBlock"},(0,l.yg)("p",null,(0,l.yg)("inlineCode",{parentName:"p"},"$ shield user edit e9fba4af-ab23-4631-abba-597b1c8e6608 --file=user.yaml"))))))}g.isMDXComponent=!0}}]);