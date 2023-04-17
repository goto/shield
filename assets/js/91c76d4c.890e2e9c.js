"use strict";(self.webpackChunkshield=self.webpackChunkshield||[]).push([[960],{3905:(e,l,i)=>{i.d(l,{Zo:()=>h,kt:()=>g});var t=i(7294);function a(e,l,i){return l in e?Object.defineProperty(e,l,{value:i,enumerable:!0,configurable:!0,writable:!0}):e[l]=i,e}function s(e,l){var i=Object.keys(e);if(Object.getOwnPropertySymbols){var t=Object.getOwnPropertySymbols(e);l&&(t=t.filter((function(l){return Object.getOwnPropertyDescriptor(e,l).enumerable}))),i.push.apply(i,t)}return i}function r(e){for(var l=1;l<arguments.length;l++){var i=null!=arguments[l]?arguments[l]:{};l%2?s(Object(i),!0).forEach((function(l){a(e,l,i[l])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(i)):s(Object(i)).forEach((function(l){Object.defineProperty(e,l,Object.getOwnPropertyDescriptor(i,l))}))}return e}function n(e,l){if(null==e)return{};var i,t,a=function(e,l){if(null==e)return{};var i,t,a={},s=Object.keys(e);for(t=0;t<s.length;t++)i=s[t],l.indexOf(i)>=0||(a[i]=e[i]);return a}(e,l);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(t=0;t<s.length;t++)i=s[t],l.indexOf(i)>=0||Object.prototype.propertyIsEnumerable.call(e,i)&&(a[i]=e[i])}return a}var d=t.createContext({}),o=function(e){var l=t.useContext(d),i=l;return e&&(i="function"==typeof e?e(l):r(r({},l),e)),i},h=function(e){var l=o(e.components);return t.createElement(d.Provider,{value:l},e.children)},p="mdxType",c={inlineCode:"code",wrapper:function(e){var l=e.children;return t.createElement(t.Fragment,{},l)}},u=t.forwardRef((function(e,l){var i=e.components,a=e.mdxType,s=e.originalType,d=e.parentName,h=n(e,["components","mdxType","originalType","parentName"]),p=o(i),u=a,g=p["".concat(d,".").concat(u)]||p[u]||c[u]||s;return i?t.createElement(g,r(r({ref:l},h),{},{components:i})):t.createElement(g,r({ref:l},h))}));function g(e,l){var i=arguments,a=l&&l.mdxType;if("string"==typeof e||a){var s=i.length,r=new Array(s);r[0]=u;var n={};for(var d in l)hasOwnProperty.call(l,d)&&(n[d]=l[d]);n.originalType=e,n[p]="string"==typeof e?e:a,r[1]=n;for(var o=2;o<s;o++)r[o]=i[o];return t.createElement.apply(null,r)}return t.createElement.apply(null,i)}u.displayName="MDXCreateElement"},2045:(e,l,i)=>{i.r(l),i.d(l,{assets:()=>d,contentTitle:()=>r,default:()=>p,frontMatter:()=>s,metadata:()=>n,toc:()=>o});var t=i(7462),a=(i(7294),i(3905));const s={},r="CLI",n={unversionedId:"reference/cli",id:"reference/cli",title:"CLI",description:"shield action",source:"@site/docs/reference/cli.md",sourceDirName:"reference",slug:"/reference/cli",permalink:"/shield/reference/cli",draft:!1,editUrl:"https://github.com/goto/shield/edit/master/docs/docs/reference/cli.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Shield API",permalink:"/shield/reference/api"}},d={},o=[{value:"shield action",id:"shield-action",level:2},{value:"shield action create flags",id:"shield-action-create-flags",level:3},{value:"shield action edit flags",id:"shield-action-edit-flags",level:3},{value:"shield action list",id:"shield-action-list",level:3},{value:"shield action view",id:"shield-action-view",level:3},{value:"shield auth",id:"shield-auth",level:2},{value:"shield completion bash|zsh|fish|powershell",id:"shield-completion-bashzshfishpowershell",level:2},{value:"shield config",id:"shield-config",level:2},{value:"shield config init",id:"shield-config-init",level:3},{value:"shield config list",id:"shield-config-list",level:3},{value:"shield environment",id:"shield-environment",level:2},{value:"shield group",id:"shield-group",level:2},{value:"shield group create flags",id:"shield-group-create-flags",level:3},{value:"shield group edit flags",id:"shield-group-edit-flags",level:3},{value:"shield group list",id:"shield-group-list",level:3},{value:"shield group view flags",id:"shield-group-view-flags",level:3},{value:"shield namespace",id:"shield-namespace",level:2},{value:"shield namespace create flags",id:"shield-namespace-create-flags",level:3},{value:"shield namespace edit flags",id:"shield-namespace-edit-flags",level:3},{value:"shield namespace list",id:"shield-namespace-list",level:3},{value:"shield namespace view",id:"shield-namespace-view",level:3},{value:"shield organization",id:"shield-organization",level:2},{value:"shield organization admadd flags",id:"shield-organization-admadd-flags",level:3},{value:"shield organization admlist",id:"shield-organization-admlist",level:3},{value:"shield organization admremove flags",id:"shield-organization-admremove-flags",level:3},{value:"shield organization create flags",id:"shield-organization-create-flags",level:3},{value:"shield organization edit flags",id:"shield-organization-edit-flags",level:3},{value:"shield organization list",id:"shield-organization-list",level:3},{value:"shield organization view flags",id:"shield-organization-view-flags",level:3},{value:"shield policy",id:"shield-policy",level:2},{value:"shield policy create flags",id:"shield-policy-create-flags",level:3},{value:"shield policy edit flags",id:"shield-policy-edit-flags",level:3},{value:"shield policy list",id:"shield-policy-list",level:3},{value:"shield policy view",id:"shield-policy-view",level:3},{value:"shield project",id:"shield-project",level:2},{value:"shield project create flags",id:"shield-project-create-flags",level:3},{value:"shield project edit flags",id:"shield-project-edit-flags",level:3},{value:"shield project list",id:"shield-project-list",level:3},{value:"shield project view flags",id:"shield-project-view-flags",level:3},{value:"shield role",id:"shield-role",level:2},{value:"shield role create flags",id:"shield-role-create-flags",level:3},{value:"shield role edit flags",id:"shield-role-edit-flags",level:3},{value:"shield role list",id:"shield-role-list",level:3},{value:"shield role view flags",id:"shield-role-view-flags",level:3},{value:"shield server",id:"shield-server",level:2},{value:"shield server init flags",id:"shield-server-init-flags",level:3},{value:"shield server migrate flags",id:"shield-server-migrate-flags",level:3},{value:"shield server migration-rollback flags",id:"shield-server-migration-rollback-flags",level:3},{value:"shield server start flags",id:"shield-server-start-flags",level:3},{value:"shield user",id:"shield-user",level:2},{value:"shield user create flags",id:"shield-user-create-flags",level:3},{value:"shield user edit flags",id:"shield-user-edit-flags",level:3},{value:"shield user list",id:"shield-user-list",level:3},{value:"shield user view flags",id:"shield-user-view-flags",level:3}],h={toc:o};function p(e){let{components:l,...i}=e;return(0,a.kt)("wrapper",(0,t.Z)({},h,i,{components:l,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"cli"},"CLI"),(0,a.kt)("h2",{id:"shield-action"},"shield action"),(0,a.kt)("p",null,"Manage actions"),(0,a.kt)("h3",{id:"shield-action-create-flags"},"shield action create ","[flags]"),(0,a.kt)("p",null,"Create an action"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string     Path to the action body file\n-H, --header string   Header <key>:<value>\n")),(0,a.kt)("h3",{id:"shield-action-edit-flags"},"shield action edit ","[flags]"),(0,a.kt)("p",null,"Edit an action"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the action body file\n")),(0,a.kt)("h3",{id:"shield-action-list"},"shield action list"),(0,a.kt)("p",null,"List all actions"),(0,a.kt)("h3",{id:"shield-action-view"},"shield action view"),(0,a.kt)("p",null,"View an action"),(0,a.kt)("h2",{id:"shield-auth"},"shield auth"),(0,a.kt)("p",null,"Auth configs that need to be used with shield"),(0,a.kt)("h2",{id:"shield-completion-bashzshfishpowershell"},"shield completion ","[bash|zsh|fish|powershell]"),(0,a.kt)("p",null,"Generate shell completion scripts"),(0,a.kt)("h2",{id:"shield-config"},"shield config"),(0,a.kt)("p",null,"Manage client configurations"),(0,a.kt)("h3",{id:"shield-config-init"},"shield config init"),(0,a.kt)("p",null,"Initialize a new client configuration"),(0,a.kt)("h3",{id:"shield-config-list"},"shield config list"),(0,a.kt)("p",null,"List client configuration settings"),(0,a.kt)("h2",{id:"shield-environment"},"shield environment"),(0,a.kt)("p",null,"List of supported environment variables"),(0,a.kt)("h2",{id:"shield-group"},"shield group"),(0,a.kt)("p",null,"Manage groups"),(0,a.kt)("h3",{id:"shield-group-create-flags"},"shield group create ","[flags]"),(0,a.kt)("p",null,"Create a group"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string     Path to the group body file\n-H, --header string   Header <key>:<value>\n")),(0,a.kt)("h3",{id:"shield-group-edit-flags"},"shield group edit ","[flags]"),(0,a.kt)("p",null,"Edit a group"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the group body file\n")),(0,a.kt)("h3",{id:"shield-group-list"},"shield group list"),(0,a.kt)("p",null,"List all groups"),(0,a.kt)("h3",{id:"shield-group-view-flags"},"shield group view ","[flags]"),(0,a.kt)("p",null,"View a group"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-m, --metadata   Set this flag to see metadata\n")),(0,a.kt)("h2",{id:"shield-namespace"},"shield namespace"),(0,a.kt)("p",null,"Manage namespaces"),(0,a.kt)("h3",{id:"shield-namespace-create-flags"},"shield namespace create ","[flags]"),(0,a.kt)("p",null,"Create a namespace"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the namespace body file\n")),(0,a.kt)("h3",{id:"shield-namespace-edit-flags"},"shield namespace edit ","[flags]"),(0,a.kt)("p",null,"Edit a namespace"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the namespace body file\n")),(0,a.kt)("h3",{id:"shield-namespace-list"},"shield namespace list"),(0,a.kt)("p",null,"List all namespaces"),(0,a.kt)("h3",{id:"shield-namespace-view"},"shield namespace view"),(0,a.kt)("p",null,"View a namespace"),(0,a.kt)("h2",{id:"shield-organization"},"shield organization"),(0,a.kt)("p",null,"Manage organizations"),(0,a.kt)("h3",{id:"shield-organization-admadd-flags"},"shield organization admadd ","[flags]"),(0,a.kt)("p",null,"add admins to an organization"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the provider config\n")),(0,a.kt)("h3",{id:"shield-organization-admlist"},"shield organization admlist"),(0,a.kt)("p",null,"list admins of an organization"),(0,a.kt)("h3",{id:"shield-organization-admremove-flags"},"shield organization admremove ","[flags]"),(0,a.kt)("p",null,"remove admins from an organization"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-u, --user string   Id of the user to be removed\n")),(0,a.kt)("h3",{id:"shield-organization-create-flags"},"shield organization create ","[flags]"),(0,a.kt)("p",null,"Create an organization"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string     Path to the organization body file\n-H, --header string   Header <key>:<value>\n")),(0,a.kt)("h3",{id:"shield-organization-edit-flags"},"shield organization edit ","[flags]"),(0,a.kt)("p",null,"Edit an organization"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the organization body file\n")),(0,a.kt)("h3",{id:"shield-organization-list"},"shield organization list"),(0,a.kt)("p",null,"List all organizations"),(0,a.kt)("h3",{id:"shield-organization-view-flags"},"shield organization view ","[flags]"),(0,a.kt)("p",null,"View an organization"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-m, --metadata   Set this flag to see metadata\n")),(0,a.kt)("h2",{id:"shield-policy"},"shield policy"),(0,a.kt)("p",null,"Manage policies"),(0,a.kt)("h3",{id:"shield-policy-create-flags"},"shield policy create ","[flags]"),(0,a.kt)("p",null,"Create a policy"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string     Path to the policy body file\n-H, --header string   Header <key>:<value>\n")),(0,a.kt)("h3",{id:"shield-policy-edit-flags"},"shield policy edit ","[flags]"),(0,a.kt)("p",null,"Edit a policy"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the policy body file\n")),(0,a.kt)("h3",{id:"shield-policy-list"},"shield policy list"),(0,a.kt)("p",null,"List all policies"),(0,a.kt)("h3",{id:"shield-policy-view"},"shield policy view"),(0,a.kt)("p",null,"View a policy"),(0,a.kt)("h2",{id:"shield-project"},"shield project"),(0,a.kt)("p",null,"Manage projects"),(0,a.kt)("h3",{id:"shield-project-create-flags"},"shield project create ","[flags]"),(0,a.kt)("p",null,"Create a project"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string     Path to the project body file\n-H, --header string   Header <key>:<value>\n")),(0,a.kt)("h3",{id:"shield-project-edit-flags"},"shield project edit ","[flags]"),(0,a.kt)("p",null,"Edit a project"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the project body file\n")),(0,a.kt)("h3",{id:"shield-project-list"},"shield project list"),(0,a.kt)("p",null,"List all projects"),(0,a.kt)("h3",{id:"shield-project-view-flags"},"shield project view ","[flags]"),(0,a.kt)("p",null,"View a project"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-m, --metadata   Set this flag to see metadata\n")),(0,a.kt)("h2",{id:"shield-role"},"shield role"),(0,a.kt)("p",null,"Manage roles"),(0,a.kt)("h3",{id:"shield-role-create-flags"},"shield role create ","[flags]"),(0,a.kt)("p",null,"Create a role"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string     Path to the role body file\n-H, --header string   Header <key>:<value>\n")),(0,a.kt)("h3",{id:"shield-role-edit-flags"},"shield role edit ","[flags]"),(0,a.kt)("p",null,"Edit a role"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the role body file\n")),(0,a.kt)("h3",{id:"shield-role-list"},"shield role list"),(0,a.kt)("p",null,"List all roles"),(0,a.kt)("h3",{id:"shield-role-view-flags"},"shield role view ","[flags]"),(0,a.kt)("p",null,"View a role"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-m, --metadata   Set this flag to see metadata\n")),(0,a.kt)("h2",{id:"shield-server"},"shield server"),(0,a.kt)("p",null,"Server management"),(0,a.kt)("h3",{id:"shield-server-init-flags"},"shield server init ","[flags]"),(0,a.kt)("p",null,"Initialize server"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},'-o, --output string      Output config file path (default "./config.yaml")\n-r, --resources string   URL path of resources. Full path prefixed with scheme where resources config yaml files are kept\n                         e.g.:\n                         local storage file "file:///tmp/resources_config"\n                         GCS Bucket "gs://shield-bucket-example"\n                         (default: file://{pwd}/resources_config)\n                         \n-u, --rule string        URL path of rules. Full path prefixed with scheme where ruleset yaml files are kept\n                         e.g.:\n                         local storage file "file:///tmp/rules"\n                         GCS Bucket "gs://shield-bucket-example"\n                         (default: file://{pwd}/rules)\n')),(0,a.kt)("h3",{id:"shield-server-migrate-flags"},"shield server migrate ","[flags]"),(0,a.kt)("p",null,"Run DB Schema Migrations"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-c, --config string   Config file path\n")),(0,a.kt)("h3",{id:"shield-server-migration-rollback-flags"},"shield server migration-rollback ","[flags]"),(0,a.kt)("p",null,"Run DB Schema Migrations Rollback to last state"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-c, --config string   Config file path\n")),(0,a.kt)("h3",{id:"shield-server-start-flags"},"shield server start ","[flags]"),(0,a.kt)("p",null,"Start server and proxy default on port 8080"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-c, --config string   Config file path\n")),(0,a.kt)("h2",{id:"shield-user"},"shield user"),(0,a.kt)("p",null,"Manage users"),(0,a.kt)("h3",{id:"shield-user-create-flags"},"shield user create ","[flags]"),(0,a.kt)("p",null,"Create an user"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string     Path to the user body file\n-H, --header string   Header <key>:<value>\n")),(0,a.kt)("h3",{id:"shield-user-edit-flags"},"shield user edit ","[flags]"),(0,a.kt)("p",null,"Edit an user"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-f, --file string   Path to the user body file\n")),(0,a.kt)("h3",{id:"shield-user-list"},"shield user list"),(0,a.kt)("p",null,"List all users"),(0,a.kt)("h3",{id:"shield-user-view-flags"},"shield user view ","[flags]"),(0,a.kt)("p",null,"View an user"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"-m, --metadata   Set this flag to see metadata\n")))}p.isMDXComponent=!0}}]);