!function(e,t){"object"==typeof exports&&"object"==typeof module?module.exports=t():"function"==typeof define&&define.amd?define([],t):"object"==typeof exports?exports.modulea=t():e.modulea=t()}(window,(function(){return function(e){var t={};function n(r){if(t[r])return t[r].exports;var o=t[r]={i:r,l:!1,exports:{}};return e[r].call(o.exports,o,o.exports,n),o.l=!0,o.exports}return n.m=e,n.c=t,n.d=function(e,t,r){n.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:r})},n.r=function(e){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},n.t=function(e,t){if(1&t&&(e=n(e)),8&t)return e;if(4&t&&"object"==typeof e&&e&&e.__esModule)return e;var r=Object.create(null);if(n.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var o in e)n.d(r,o,function(t){return e[t]}.bind(null,o));return r},n.n=function(e){var t=e&&e.__esModule?function(){return e.default}:function(){return e};return n.d(t,"a",t),t},n.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},n.p="",n(n.s=0)}([function(e,t,n){"use strict";n.r(t),n.d(t,"Component",(function(){return r}));let r=n(1)},function(e,t,n){"use strict";n.r(t);var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticStyle:{border:"2px solid red",margin:"5px",padding:"5px",width:"500px",height:"400px"}},[n("div",[n("div",[e._v("\n            RTSP:"),n("input",{directives:[{name:"model",rawName:"v-model",value:e.rtspUrl,expression:"rtspUrl"}],domProps:{value:e.rtspUrl},on:{input:function(t){t.target.composing||(e.rtspUrl=t.target.value)}}})])]),e._v(" "),n("button",{attrs:{type:"button"},on:{click:e.test}},[e._v("连接webrtc")]),e._v(" "),e._m(0)])};r._withStripped=!0;var o={iceServers:[{url:"stun:stun.l.google.com:19302"}]};var i=function(e,t,n,r,o,i,a,s){var c,d="function"==typeof e?e.options:e;if(t&&(d.render=t,d.staticRenderFns=n,d._compiled=!0),r&&(d.functional=!0),i&&(d._scopeId="data-v-"+i),a?(c=function(e){(e=e||this.$vnode&&this.$vnode.ssrContext||this.parent&&this.parent.$vnode&&this.parent.$vnode.ssrContext)||"undefined"==typeof __VUE_SSR_CONTEXT__||(e=__VUE_SSR_CONTEXT__),o&&o.call(this,e),e&&e._registeredComponents&&e._registeredComponents.add(a)},d._ssrRegister=c):o&&(c=s?function(){o.call(this,this.$root.$options.shadowRoot)}:o),c)if(d.functional){d._injectStyles=c;var u=d.render;d.render=function(e,t){return c.call(t),u(e,t)}}else{var l=d.beforeCreate;d.beforeCreate=l?[].concat(l,c):[c]}return{exports:e,options:d}}({data:()=>({rtspUrl:"rtsp://admin:DFwl123456@172.16.133.207:554/Streaming/Channels/101?transportmode=unicast&profile=Profile_1"}),methods:{test(){var e=this;navigator.getUserMedia=navigator.getUserMedia||navigator.webkitGetUserMedia||navigator.mozGetUserMedia;var t=new RTCPeerConnection(o);t.onnegotiationneeded=function(){t.createOffer((function(n){t.setLocalDescription(n),$.ajax({url:"/api/recive",type:"POST",data:{url:e.rtspUrl,data:Base64.encode(n.sdp)}}).then(e=>{t.setRemoteDescription(new RTCSessionDescription({type:"answer",sdp:atob(e)}))})}),(function(e,t,n){}))};let n=e=>{console.log(e)};t.ontrack=function(e){n(e.streams.length+" track is delivered");var t=document.getElementById("video");t.srcObject=e.streams[0],t.muted=!0,t.autoplay=!0,t.controls=!0},t.oniceconnectionstatechange=e=>n(t.iceConnectionState),t.addTransceiver("video",{direction:"sendrecv"})}}},r,[function(){var e=this.$createElement,t=this._self._c||e;return t("div",[t("video",{attrs:{id:"video",width:"500","height:300":""}})])}],!1,null,null,null);i.options.__file="src/moduleb/index.vue";t.default=i.exports}])}));