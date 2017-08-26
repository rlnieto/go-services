eventos.controller("eventosCtrl", ["$scope", "$resource", "eventosFactory", function($scope, $resource, eventosFactory){
    $scope.error = false;

    // Recupera los eventos asociados a un usuario
    eventosFactory.eventosDelUsuario(2)
      .then(function(data){
        $scope.eventos = data;
      })
      .catch(function(err){
        $scope.error = true;
        $scope.mensajeError = err.data.Error;
      });

}]);
