"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[465],{5162:(e,a,t)=>{t.d(a,{Z:()=>r});var n=t(7294),l=t(6010);const o="tabItem_Ymn6";function r(e){let{children:a,hidden:t,className:r}=e;return n.createElement("div",{role:"tabpanel",className:(0,l.Z)(o,r),hidden:t},a)}},5488:(e,a,t)=>{t.d(a,{Z:()=>m});var n=t(7462),l=t(7294),o=t(6010),r=t(2389),s=t(7392),i=t(7094),c=t(2466);const u="tabList__CuJ",d="tabItem_LNqP";function p(e){const{lazy:a,block:t,defaultValue:r,values:p,groupId:m,className:b}=e,f=l.Children.map(e.children,(e=>{if((0,l.isValidElement)(e)&&"value"in e.props)return e;throw new Error(`Docusaurus error: Bad <Tabs> child <${"string"==typeof e.type?e.type:e.type.name}>: all children of the <Tabs> component should be <TabItem>, and every <TabItem> should have a unique "value" prop.`)})),h=p??f.map((e=>{let{props:{value:a,label:t,attributes:n}}=e;return{value:a,label:t,attributes:n}})),g=(0,s.l)(h,((e,a)=>e.value===a.value));if(g.length>0)throw new Error(`Docusaurus error: Duplicate values "${g.map((e=>e.value)).join(", ")}" found in <Tabs>. Every value needs to be unique.`);const T=null===r?r:r??f.find((e=>e.props.default))?.props.value??f[0].props.value;if(null!==T&&!h.some((e=>e.value===T)))throw new Error(`Docusaurus error: The <Tabs> has a defaultValue "${T}" but none of its children has the corresponding value. Available values are: ${h.map((e=>e.value)).join(", ")}. If you intend to show no default tab, use defaultValue={null} instead.`);const{tabGroupChoices:v,setTabGroupChoices:k}=(0,i.U)(),[y,I]=(0,l.useState)(T),N=[],{blockElementScrollPositionUntilNextRender:Z}=(0,c.o5)();if(null!=m){const e=v[m];null!=e&&e!==y&&h.some((a=>a.value===e))&&I(e)}const E=e=>{const a=e.currentTarget,t=N.indexOf(a),n=h[t].value;n!==y&&(Z(a),I(n),null!=m&&k(m,String(n)))},x=e=>{let a=null;switch(e.key){case"Enter":E(e);break;case"ArrowRight":{const t=N.indexOf(e.currentTarget)+1;a=N[t]??N[0];break}case"ArrowLeft":{const t=N.indexOf(e.currentTarget)-1;a=N[t]??N[N.length-1];break}}a?.focus()};return l.createElement("div",{className:(0,o.Z)("tabs-container",u)},l.createElement("ul",{role:"tablist","aria-orientation":"horizontal",className:(0,o.Z)("tabs",{"tabs--block":t},b)},h.map((e=>{let{value:a,label:t,attributes:r}=e;return l.createElement("li",(0,n.Z)({role:"tab",tabIndex:y===a?0:-1,"aria-selected":y===a,key:a,ref:e=>N.push(e),onKeyDown:x,onClick:E},r,{className:(0,o.Z)("tabs__item",d,r?.className,{"tabs__item--active":y===a})}),t??a)}))),a?(0,l.cloneElement)(f.filter((e=>e.props.value===y))[0],{className:"margin-top--md"}):l.createElement("div",{className:"margin-top--md"},f.map(((e,a)=>(0,l.cloneElement)(e,{key:a,hidden:e.props.value!==y})))))}function m(e){const a=(0,r.Z)();return l.createElement(p,(0,n.Z)({key:String(a)},e))}},9380:(e,a,t)=>{t.r(a),t.d(a,{assets:()=>d,contentTitle:()=>c,default:()=>b,frontMatter:()=>i,metadata:()=>u,toc:()=>p});var n=t(7462),l=(t(7294),t(3905)),o=t(5488),r=t(5162),s=t(814);const i={},c="Managing Relations",u={unversionedId:"guides/managing-relation",id:"guides/managing-relation",title:"Managing Relations",description:"A relation in Shield looks like",source:"@site/docs/guides/managing-relation.md",sourceDirName:"guides",slug:"/guides/managing-relation",permalink:"/shield/guides/managing-relation",draft:!1,editUrl:"https://github.com/odpf/shield/edit/master/docs/docs/guides/managing-relation.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Manage Resources",permalink:"/shield/guides/managing-resource"},next:{title:"Managing Users",permalink:"/shield/guides/managing-user"}},d={},p=[{value:"API Interface",id:"api-interface",level:2},{value:"Create Relations",id:"create-relations",level:3},{value:"List Relations",id:"list-relations",level:3},{value:"Get Relations",id:"get-relations",level:3},{value:"Delete relation",id:"delete-relation",level:3}],m={toc:p};function b(e){let{components:a,...t}=e;return(0,l.kt)("wrapper",(0,n.Z)({},m,t,{components:a,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"managing-relations"},"Managing Relations"),(0,l.kt)("p",null,"A relation in Shield looks like"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "relations": [\n        {\n            "id": "08effbce-42cb-4b7e-a808-ad17cd3445df",\n            "objectId": "a9f784cf-0f29-486f-92d0-51300295f7e8",\n            "objectNamespace": "entropy/firehose",\n            "subject": "user:598688c6-8c6d-487f-b324-ef3f4af120bb",\n            "roleName": "entropy/firehose:owner",\n            "createdAt": null,\n            "updatedAt": null\n        }\n    ]\n}\n')),(0,l.kt)("h2",{id:"api-interface"},"API Interface"),(0,l.kt)("h3",{id:"create-relations"},"Create Relations"),(0,l.kt)(o.Z,{groupId:"api",mdxType:"Tabs"},(0,l.kt)(r.Z,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.kt)(s.Z,{className:"language-bash",mdxType:"CodeBlock"},'$ curl --location --request POST \'http://localhost:8000/admin/v1beta1/relations\'\n--header \'Content-Type: application/json\'\n--header \'Accept: application/json\'\n--data-raw \'{\n  "objectId": "a9f784cf-0f29-486f-92d0-51300295f7e8",\n  "objectNamespace": "entropy/firehose",\n  "subject": "user:doe.john@odpf.io",\n  "roleName": "owner"\n}\''))),(0,l.kt)("h3",{id:"list-relations"},"List Relations"),(0,l.kt)(o.Z,{groupId:"api",mdxType:"Tabs"},(0,l.kt)(r.Z,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.kt)(s.Z,{className:"language-bash",mdxType:"CodeBlock"},"$ curl --location --request GET 'http://localhost:8000/admin/v1beta1/relations'\n--header 'Accept: application/json'"))),(0,l.kt)("h3",{id:"get-relations"},"Get Relations"),(0,l.kt)(o.Z,{groupId:"api",mdxType:"Tabs"},(0,l.kt)(r.Z,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.kt)(s.Z,{className:"language-bash",mdxType:"CodeBlock"},"$ curl --location --request GET 'http://localhost:8000/admin/v1beta1/relations/f959a605-8755-4ee4-b898-a1e26f596c4d'\n--header 'Accept: application/json'"))),(0,l.kt)("h3",{id:"delete-relation"},"Delete relation"),(0,l.kt)(o.Z,{groupId:"api",mdxType:"Tabs"},(0,l.kt)(r.Z,{value:"HTTP",label:"HTTP",default:!0,mdxType:"TabItem"},(0,l.kt)(s.Z,{className:"language-bash",mdxType:"CodeBlock"},"$ curl --location --request DELETE 'http://localhost:8000/admin/v1beta1/\n    object/a9f784cf-0f29-486f-92d0-51300295f7e8/\n    subject/448d52d4-48cb-495e-8ec5-8afc55c624ca/\n    role/owner'\n--header 'Accept: application/json'"))))}b.isMDXComponent=!0}}]);