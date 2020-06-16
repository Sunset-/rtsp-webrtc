!(function (e, t) {
  "object" == typeof exports && "object" == typeof module
    ? (module.exports = t())
    : "function" == typeof define && define.amd
    ? define([], t)
    : "object" == typeof exports
    ? (exports.modulea = t())
    : (e.modulea = t());
})(window, function () {
  return (function (e) {
    var t = {};
    function n(o) {
      if (t[o]) return t[o].exports;
      var r = (t[o] = { i: o, l: !1, exports: {} });
      return e[o].call(r.exports, r, r.exports, n), (r.l = !0), r.exports;
    }
    return (
      (n.m = e),
      (n.c = t),
      (n.d = function (e, t, o) {
        n.o(e, t) || Object.defineProperty(e, t, { enumerable: !0, get: o });
      }),
      (n.r = function (e) {
        "undefined" != typeof Symbol &&
          Symbol.toStringTag &&
          Object.defineProperty(e, Symbol.toStringTag, { value: "Module" }),
          Object.defineProperty(e, "__esModule", { value: !0 });
      }),
      (n.t = function (e, t) {
        if ((1 & t && (e = n(e)), 8 & t)) return e;
        if (4 & t && "object" == typeof e && e && e.__esModule) return e;
        var o = Object.create(null);
        if (
          (n.r(o),
          Object.defineProperty(o, "default", { enumerable: !0, value: e }),
          2 & t && "string" != typeof e)
        )
          for (var r in e)
            n.d(
              o,
              r,
              function (t) {
                return e[t];
              }.bind(null, r)
            );
        return o;
      }),
      (n.n = function (e) {
        var t =
          e && e.__esModule
            ? function () {
                return e.default;
              }
            : function () {
                return e;
              };
        return n.d(t, "a", t), t;
      }),
      (n.o = function (e, t) {
        return Object.prototype.hasOwnProperty.call(e, t);
      }),
      (n.p = ""),
      n((n.s = 0))
    );
  })([
    function (e, t, n) {
      "use strict";
      n.r(t),
        n.d(t, "Init", function () {
          return i;
        }),
        n.d(t, "Component", function () {
          return u;
        });
      var o = function () {
        var e = this.$createElement;
        return (this._self._c || e)("div", [this._v("我是一个VUE组件-MODULEA")]);
      };
      o._withStripped = !0;
      var r = (function (e, t, n, o, r, i, u, f) {
        var s,
          c = "function" == typeof e ? e.options : e;
        if (
          (t && ((c.render = t), (c.staticRenderFns = n), (c._compiled = !0)),
          o && (c.functional = !0),
          i && (c._scopeId = "data-v-" + i),
          u
            ? ((s = function (e) {
                (e =
                  e ||
                  (this.$vnode && this.$vnode.ssrContext) ||
                  (this.parent &&
                    this.parent.$vnode &&
                    this.parent.$vnode.ssrContext)) ||
                  "undefined" == typeof __VUE_SSR_CONTEXT__ ||
                  (e = __VUE_SSR_CONTEXT__),
                  r && r.call(this, e),
                  e &&
                    e._registeredComponents &&
                    e._registeredComponents.add(u);
              }),
              (c._ssrRegister = s))
            : r &&
              (s = f
                ? function () {
                    r.call(this, this.$root.$options.shadowRoot);
                  }
                : r),
          s)
        )
          if (c.functional) {
            c._injectStyles = s;
            var l = c.render;
            c.render = function (e, t) {
              return s.call(t), l(e, t);
            };
          } else {
            var d = c.beforeCreate;
            c.beforeCreate = d ? [].concat(d, s) : [s];
          }
        return { exports: e, options: c };
      })({}, o, [], !1, null, null, null);
      function i() {
        console.log("模块A初始化了");
      }
      r.options.__file = "index.vue";
      let u = r.exports;
    },
  ]);
});
