function generateRoute(item) {
  var r = {
    path: item.root?"/":("/" + item.module),
    name: item.module,
    redirect :item.redirect,
    component: () => {
      return $module.loadModuleComponent(item.module, item.version);
    },
  };
  if (item.children) {
    r.children = item.children.map((item) => {
      return generateRoute(item);
    });
  }
  return r;
}

$module.loadMenu().then((res) => {
  //路由可来自接口或配置文件
  var routes = [];
  if (res) {
    Object.keys(res).forEach((k) => {
      routes.push(generateRoute(res[k]));
    });
  }
  const router = new VueRouter();
  router.addRoutes(routes);
  window.$router = router;


  //启动
  const app = new Vue({
    router,
  }).$mount("#app");
});
