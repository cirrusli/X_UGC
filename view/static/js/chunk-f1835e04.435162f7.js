(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-f1835e04"],{1662:function(e,t,n){"use strict";n("4401")},"408a":function(e,t,n){var r=n("e330");e.exports=r(1..valueOf)},4401:function(e,t,n){},"4de4":function(e,t,n){"use strict";var r=n("23e7"),s=n("b727").filter,o=n("1dde"),a=o("filter");r({target:"Array",proto:!0,forced:!a},{filter:function(e){return s(this,e,arguments.length>1?arguments[1]:void 0)}})},5530:function(e,t,n){"use strict";n.d(t,"a",(function(){return o}));n("b64b"),n("a4d3"),n("4de4"),n("d3b7"),n("e439"),n("159b"),n("dbb4");var r=n("ade3");function s(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?s(Object(n),!0).forEach((function(t){Object(r["a"])(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):s(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}},5899:function(e,t){e.exports="\t\n\v\f\r                　\u2028\u2029\ufeff"},"58a8":function(e,t,n){var r=n("e330"),s=n("1d80"),o=n("577e"),a=n("5899"),i=r("".replace),c="["+a+"]",u=RegExp("^"+c+c+"*"),f=RegExp(c+c+"*$"),d=function(e){return function(t){var n=o(s(t));return 1&e&&(n=i(n,u,"")),2&e&&(n=i(n,f,"")),n}};e.exports={start:d(1),end:d(2),trim:d(3)}},"71b2":function(e,t,n){"use strict";n.r(t);var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[n("van-nav-bar",{attrs:{title:e.userInfo.username,fixed:"","left-text":"返回","left-arrow":""},on:{"click-left":e.onClickLeft}}),n("div",{ref:"chattingRecords",staticClass:"chattingRecords"},[n("div",{ref:"msgBox"},e._l(e.messageList,(function(t,r){return n("div",{key:r,class:"send"==t[0]?"message2 messageItem":"message1 messageItem"},[n("van-image",{attrs:{round:"",width:"40px",height:"40px",src:"send"==t[0]?e.proFilePhoto:e.userInfo.head_photo}}),n("p",[e._v(" "+e._s(t[1])+" ")])],1)})),0)]),n("div",{staticClass:"inputbox"},[n("div",{staticClass:"saybox"},[n("van-icon",{attrs:{name:"volume-o"}})],1),n("van-field",{staticClass:"inputdiv",attrs:{rows:"1",autosize:"",type:"textarea",placeholder:"请输入..."},model:{value:e.sendInfo.content,callback:function(t){e.$set(e.sendInfo,"content",t)},expression:"sendInfo.content"}}),n("van-button",{attrs:{color:"linear-gradient(to right, #ff6034, #ee0a24)",size:"small"},on:{click:e.sendMessage}},[e._v(" 发送 ")])],1)],1)},s=[],o=n("1da1"),a=n("5530"),i=(n("96cf"),n("a9e3"),n("d3b7"),n("159b"),n("ac1f"),n("1276"),n("d9e2"),n("c1df")),c=n.n(i),u=n("fa7d"),f={data:function(){return{sendInfo:{type:1,content:"",send_time:c()().format("YYYY-MM-DD HH:mm:ss.SSS"),client_ids:[]},userId:"",messageList:[],userInfo:{},proFilePhoto:""}},computed:{msgarr:function(){try{return this.$store.getters["messInfo/getMsg"](Number(this.$route.query.userId))}catch(e){console.log("computed",e)}}},watch:{msgarr:{handler:function(e,t){var n=this;try{if(null!=e){this.userInfo=Object(u["a"])(e.user_info);var r=e.history_data;this.messageList=[],r.forEach((function(e){var t=e.split("+");n.messageList.push([t[1],t[2]])}))}}catch(s){console.log(s)}},deep:!0,immediate:!0}},methods:{sendMessage:function(){this.sendInfo.client_ids=["".concat(this.userId)],this.$ws.sendSock(this.sendInfo),this.$store.dispatch("messInfo/addUserMsg",Object(a["a"])(Object(a["a"])({},this.sendInfo),{},{from_user:this.userId,msgType:"send"})),this.sendInfo.content=""},onClickLeft:function(){this.$router.push("/message")},getMsg:function(){var e=this;return Object(o["a"])(regeneratorRuntime.mark((function t(){var n,r;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,e.$http.get("/userInfo/getUserInfo");case 2:n=t.sent,r=n.data,"success"===r.status&&(r=Object(u["a"])(r),e.proFilePhoto=r.data.userinfo.head_photo);case 5:case"end":return t.stop()}}),t)})))()},onBottom:function(){var e=this;this.$nextTick((function(){var t=e.$refs.msgBox;e.messageList.length>5&&window.scrollTo(0,t.offsetHeight)}))}},updated:function(){this.onBottom()},mounted:function(){var e=this;this.userId=Number(this.$route.query.userId),this.getMsg(),this.$nextTick((function(){e.onBottom()}))},beforeDestroy:function(){var e=this;return Object(o["a"])(regeneratorRuntime.mark((function t(){var n,r;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,e.$http.put("/notify/resetUnreadChat?chat_userid=".concat(e.userId));case 2:if(n=t.sent,r=n.data,"success"===r.status){t.next=6;break}throw new Error("重置聊天失败");case 6:case"end":return t.stop()}}),t)})))()}},d=f,l=(n("1662"),n("2877")),h=Object(l["a"])(d,r,s,!1,null,"651884ab",null);t["default"]=h.exports},a9e3:function(e,t,n){"use strict";var r=n("83ab"),s=n("da84"),o=n("e330"),a=n("94ca"),i=n("6eeb"),c=n("1a2d"),u=n("7156"),f=n("3a9b"),d=n("d9b5"),l=n("c04e"),h=n("d039"),p=n("241c").f,b=n("06cf").f,g=n("9bf2").f,m=n("408a"),v=n("58a8").trim,I="Number",O=s[I],w=O.prototype,y=s.TypeError,x=o("".slice),j=o("".charCodeAt),_=function(e){var t=l(e,"number");return"bigint"==typeof t?t:N(t)},N=function(e){var t,n,r,s,o,a,i,c,u=l(e,"number");if(d(u))throw y("Cannot convert a Symbol value to a number");if("string"==typeof u&&u.length>2)if(u=v(u),t=j(u,0),43===t||45===t){if(n=j(u,2),88===n||120===n)return NaN}else if(48===t){switch(j(u,1)){case 66:case 98:r=2,s=49;break;case 79:case 111:r=8,s=55;break;default:return+u}for(o=x(u,2),a=o.length,i=0;i<a;i++)if(c=j(o,i),c<48||c>s)return NaN;return parseInt(o,r)}return+u};if(a(I,!O(" 0o1")||!O("0b1")||O("+0x1"))){for(var E,k=function(e){var t=arguments.length<1?0:O(_(e)),n=this;return f(w,n)&&h((function(){m(n)}))?u(Object(t),n,k):t},P=r?p(O):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,isFinite,isInteger,isNaN,isSafeInteger,parseFloat,parseInt,fromString,range".split(","),$=0;P.length>$;$++)c(O,E=P[$])&&!c(k,E)&&g(k,E,b(O,E));k.prototype=w,w.constructor=k,i(s,I,k)}},ade3:function(e,t,n){"use strict";function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}n.d(t,"a",(function(){return r}))},dbb4:function(e,t,n){var r=n("23e7"),s=n("83ab"),o=n("56ef"),a=n("fc6a"),i=n("06cf"),c=n("8418");r({target:"Object",stat:!0,sham:!s},{getOwnPropertyDescriptors:function(e){var t,n,r=a(e),s=i.f,u=o(r),f={},d=0;while(u.length>d)n=s(r,t=u[d++]),void 0!==n&&c(f,t,n);return f}})},e439:function(e,t,n){var r=n("23e7"),s=n("d039"),o=n("fc6a"),a=n("06cf").f,i=n("83ab"),c=s((function(){a(1)})),u=!i||c;r({target:"Object",stat:!0,forced:u,sham:!i},{getOwnPropertyDescriptor:function(e,t){return a(o(e),t)}})}}]);
//# sourceMappingURL=chunk-f1835e04.435162f7.js.map