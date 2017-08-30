//-----------------------------------------------------------------------------
// Listado de los eventos de un usuario
//
//-----------------------------------------------------------------------------
eventos.controller("listaEventosCtrl", ["$scope", "$resource", "eventosFactory", function($scope, $resource, eventosFactory){
    $scope.asistencia = "";
    $scope.mensaje = "";
    $scope.error = false;

    //-----------------------------------------------------------------------
    // Recupera los eventos asociados a un usuario
    //-----------------------------------------------------------------------
    eventosFactory.eventosDelUsuario(2)
      .then(function(data){
        $scope.eventos = data;
      })
      .catch(function(err){
        $scope.error = true;
        $scope.mensajeError = err.data.Error;
      });

      //-----------------------------------------------------------------------
      // Modificación del indicador de asistencia al evento
      //-----------------------------------------------------------------------
      $scope.modificarAsistencia = function(usuario, evento, indicador){

        eventosFactory.indicadorAsistenciaEvento(usuario, evento, indicador)
          .then(function(data){
            //alert("Modificación correcta");
            //$scope.eventos = data;
            $scope.mensaje = "Modificación correcta!";
          })
          .catch(function(err){
            $scope.error = true;
            $scope.mensaje = err.data.Error;
          });

      }


}]);


//-----------------------------------------------------------------------------
// Detalle de un evento
//
//-----------------------------------------------------------------------------
eventos.controller("detalleEventoCtrl", ["$scope", "$resource", function($scope, $resource){


}]);


//-----------------------------------------------------------------------------
// Pagina 404
//
//-----------------------------------------------------------------------------
eventos.controller("noEncontradoCtrl", ["$scope", "$resource", function($scope, $resource){


}]);
