(function (window) {
  var $module = (window.$module = {});

  var loadModules = {};

  //加载菜单
  $module.loadMenu = function () {
    return fetch(`/menu.json`).then((res) => {
      return res.json();
    });
  };

  //加载模块
  $module.loadModule = function (moduleName, version) {
    if (loadModules[moduleName]) {
      return loadModules[moduleName];
    }
    return (loadModules[moduleName] = fetch(`/modules/${moduleName}.js`)
      .then((res) => {
        return res.text();
      })
      .then((sourceCode) => {
        var m = {
          exports: {},
        };
        Function(`
            return function(module, exports, require) {
              ${sourceCode}
            }
          `)()(m, m.exports);
        return m.exports;
      }));
  };

  //加载模块
  $module.loadModuleComponent = function (moduleName, version) {
    return $module.loadModule(moduleName, version).then((m) => m.Component);
  };

  //模块调用
  $module.call = function (moduleName, funcName, args) {
    return $module.loadModule(moduleName).then((m) => {
      if (typeof m[funcName] == "function") {
        return m[funcName].apply(m, args);
      } else {
        console.error(`module [${moduleName}] not hava func [${funcName}]`);
      }
    });
  };

  //模块跳转
  $module.jump = function (moduleName, funcName, args) {
    if($router.currentRoute.name==moduleName){
      return;
    }
    return $router.push({
      name: moduleName,
    });
  };
})(window);
