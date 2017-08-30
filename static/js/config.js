eventos.config(function($routeProvider){
  $routeProvider.
  when("/listaEventos", {
    templateUrl: "html/lista_eventos.html",
    controller: "listaEventosCtrl"
  }).
  when("/detalleEvento", {
    templateUrl: "html/detalle_evento.html",
    controller: "detalleEventoCtrl"
  }).
  when("/error", {
    templateUrl: "html/404.html",
    controller: "noEncontradoCtrl"
  }).
  otherwise({
    redirectTo: '/listaEventos'
  });
});
