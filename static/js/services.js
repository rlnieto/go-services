//-----------------------------------------------------------------------------
// Llamadas a servicios relacionados con los eventos
//
//-----------------------------------------------------------------------------
eventos.factory('eventosFactory', ['$resource', '$q', function($resource, $q){
  var urlEventos = 'http://localhost:8080/api/eventosusuario';
  return{

    //-------------------------------------------------
    // Consulta de los eventos de un usuario
    //-------------------------------------------------
    eventosDelUsuario: function(usuario){
      var defered  = $q.defer();
      //var promise = defered.promise;

      var Resource = $resource(urlEventos, {idUsuario: 'idusuario'},{
        'query': {method:'GET', isArray:true}
      });

      Resource.query({idusuario:usuario},
        function(data){
          defered.resolve(data);
        },
        function(error){
          console.log(error.data.Error);
          defered.reject(error);
        });
      return defered.promise;
    }
  }
}]);
