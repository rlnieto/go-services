//-----------------------------------------------------------------------------
// Llamadas a servicios relacionados con los eventos
//
//-----------------------------------------------------------------------------
eventos.factory('eventosFactory', ['$resource', '$q', function($resource, $q){
  var urlEventos = 'http://localhost:8080/api/eventosusuario';
  var urlUsuariosEvento = 'http://localhost:8080/api/usuariosevento'

  return{

    //-------------------------------------------------
    // Consulta de los eventos de un usuario
    //-------------------------------------------------
    eventosDelUsuario: function(usuario){
      var defered  = $q.defer();
      //var promise = defered.promise;

      var Resource = $resource(urlEventos, {},{
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
    },


    //-------------------------------------------------
    // Modifica el indicador de asistencia de un
    // usuario a un evento
    //-------------------------------------------------
    indicadorAsistenciaEvento: function(usuario, evento, indicador){
      var defered  = $q.defer();
      //var promise = defered.promise;

      var Resource = $resource(urlUsuariosEvento, {idusuario: usuario, idevento: evento, confirmado: indicador},
      {
        'update': { method: 'PUT'}
      });

      Resource.update({},
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
