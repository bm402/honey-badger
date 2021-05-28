(this.webpackJsonpdashboard=this.webpackJsonpdashboard||[]).push([[0],{49:function(t,e,a){},50:function(t,e,a){},51:function(t,e,a){},52:function(t,e,a){},59:function(t,e,a){},60:function(t,e,a){},62:function(t,e,a){"use strict";a.r(e);var s=a(17),c=a.n(s),n=a(31),o=a(22),i=a(23),r=a(1),d=function(){return Object(r.jsxs)(o.a,{bg:"light",expand:"lg",children:[Object(r.jsx)(o.a.Brand,{href:"/honey-badger/",children:"Honey Badger"}),Object(r.jsx)(o.a.Toggle,{"aria-controls":"basic-navbar-nav"}),Object(r.jsxs)(o.a.Collapse,{id:"basic-navbar-nav",children:[Object(r.jsx)(i.a,{children:Object(r.jsx)(i.a.Link,{href:"/honey-badger/heatmap",children:"Heatmap"})}),Object(r.jsx)(i.a,{children:Object(r.jsx)(i.a.Link,{href:"/honey-badger/stats",children:"Stats"})})]})]})},j=a(6),l=function(){return Object(r.jsx)("div",{children:"This is the homepage"})},m=a(65),u=a(66),b=a(0),p=a.n(b),h=a(64),x=a(21),O=a.n(x),g=(a(48),function(){var t=Object(h.a)();return p.a.useEffect((function(){fetch("https://omf1aavgfc.execute-api.eu-west-2.amazonaws.com/prod/v1/heatmap-data").then((function(t){return t.json()})).then((function(t){return t.heatmap_data_points})).then((function(e){for(var a=e.map((function(t){return[t.lat,t.lon,t.count]})).sort((function(t,e){return t[2]-e[2]})),s=1,c=a[0][2],n=0;n<a.length;n++)a[n][2]>c&&s++,c=a[n][2],a[n][2]=s;var o=a[a.length-1][2],i=a.map((function(t){return[t[0],t[1],t[2]/o]}));O.a.heatLayer(i,{minOpacity:.4,maxZoom:3,radius:20,blur:15,gradient:{.4:"blue",.6:"lime",.8:"yellow",.9:"orange",1:"red"}}).addTo(t)})).catch(console.log)}),[t]),null}),f=(a(49),function(){return Object(r.jsxs)(m.a,{className:"heatmap",center:[0,0],zoom:2,children:[Object(r.jsx)(g,{}),Object(r.jsx)(u.a,{attribution:'\xa9 <a href="http://osm.org/copyright">OpenStreetMap</a> contributors',url:"https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"})]})}),v=a(18),y=a(35),_=a(8),N=a.p+"static/media/gold.7b323ec2.jpg",L=a.p+"static/media/silver.9a7c5b0e.jpg",D=a.p+"static/media/bronze.d1bf6559.jpg",w=(a(50),a(51),function(t){return Object(r.jsxs)(y.a,{className:"podium",children:[Object(r.jsx)(_.a,{className:"podium-title-card",children:Object(r.jsx)(_.a.Body,{className:"podium-title-text",children:Object(r.jsx)(_.a.Text,{children:t.title})})}),Object(r.jsxs)(_.a,{children:[Object(r.jsxs)(_.a.Body,{children:[Object(r.jsx)(_.a.Img,{className:"podium-medal-image",src:N}),Object(r.jsx)(_.a.Text,{className:"".concat(t.isDataLoaded?"":"loading"," ").concat("podium-text-"+t.type),children:t.data[0].value})]}),Object(r.jsx)(_.a.Footer,{className:"podium-gold-footer"})]}),Object(r.jsxs)(_.a,{children:[Object(r.jsxs)(_.a.Body,{children:[Object(r.jsx)(_.a.Img,{className:"podium-medal-image",src:L}),Object(r.jsx)(_.a.Text,{className:"".concat(t.isDataLoaded?"":"loading"," ").concat("podium-text-"+t.type),children:t.data[1].value})]}),Object(r.jsx)(_.a.Footer,{className:"podium-silver-footer"})]}),Object(r.jsxs)(_.a,{children:[Object(r.jsxs)(_.a.Body,{children:[Object(r.jsx)(_.a.Img,{className:"podium-medal-image",src:D}),Object(r.jsx)(_.a.Text,{className:"".concat(t.isDataLoaded?"":"loading"," ").concat("podium-text-"+t.type),children:t.data[2].value})]}),Object(r.jsx)(_.a.Footer,{className:"podium-bronze-footer"})]})]})}),B=(a(52),function(){var t=Object(b.useState)(!1),e=Object(v.a)(t,2),a=e[0],s=e[1],c=Object(b.useState)({most_connections:[{},{},{}],most_active_cities:[{},{},{}],most_active_countries:[{},{},{}],most_ip_addresses:[{},{},{}],most_ingress_ports:[{},{},{}]}),n=Object(v.a)(c,2),o=n[0],i=n[1];return p.a.useEffect((function(){fetch("https://omf1aavgfc.execute-api.eu-west-2.amazonaws.com/prod/v1/stats-data").then((function(t){return t.json()})).then((function(t){i(t),s(!0)})).catch(console.log)}),[]),Object(r.jsxs)("div",{className:"stats-page",children:[Object(r.jsx)(w,{title:"Most connection attempts",data:o.most_connections,type:"number",isDataLoaded:a}),Object(r.jsx)(w,{title:"Most active city",data:o.most_active_cities,type:"string",isDataLoaded:a}),Object(r.jsx)(w,{title:"Most active country",data:o.most_active_countries,type:"string",isDataLoaded:a}),Object(r.jsx)(w,{title:"Most IP addresses used",data:o.most_ip_addresses,type:"number",isDataLoaded:a}),Object(r.jsx)(w,{title:"Most ingress ports tried",data:o.most_ingress_ports,type:"number",isDataLoaded:a})]})}),T=function(){return Object(r.jsxs)(j.c,{children:[Object(r.jsx)(j.a,{exact:!0,path:"/honey-badger/",component:l}),Object(r.jsx)(j.a,{exact:!0,path:"/honey-badger/heatmap",component:f}),Object(r.jsx)(j.a,{exact:!0,path:"/honey-badger/stats",component:B})]})},z=(a(59),function(){return Object(r.jsxs)("div",{className:"App",children:[Object(r.jsx)(d,{}),Object(r.jsx)(T,{})]})});a(60),a(61);c.a.render(Object(r.jsx)(n.a,{children:Object(r.jsx)(z,{})}),document.getElementById("root"))}},[[62,1,2]]]);
//# sourceMappingURL=main.bb575780.chunk.js.map