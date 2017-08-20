parking.directive("alert", function(){
  return{
    restrict: 'E',   // A - atributo, E - elemento, C - clase, M - comentario
    scope:{
      topic: '=',
      description: '=',
      close: '&'
    },
    templateUrl: "alert.html",
    replace: true,
    transclude: true
  };
});

parking.directive("accordionItem", function(){
  return{
    templateUrl: "accordionItem.html",
    restrict: 'E',
    scope: {
      title: '@'
    },
    transclude: true,
    require: "^accordion",   // ? - busca el controller, ^ - controller en el padre
    link: function(scope, element, attrs, ctrl, transcludeFn){
      ctrl.addAccordionItem(scope);   // invoca una función en el controlador del padre
      element.bind("click", function(){
        scope.$apply(function(){
          ctrl.closeAll();  // Cierra todos
          scope.active = !scope.active;  // Y abre el suyo
        });
      });
    }
  };
});

parking.directive("accordion", function(){
  return{
    template: '<div ng-transclude></div>',
    restrict: "E",
    transclude: true,
      controller: function ($scope, $element, $attrs, $transclude){
        var accordionItems = [];

        var addAccordionItem = function (accordionScope){
          accordionItems.push(accordionScope);
        }

        var closeAll = function(){
          angular.forEach(accordionItems, function (accordionScope){
            accordionScope.active = false;
          });
        }

        return{
          addAccordionItem: addAccordionItem,
          closeAll: closeAll
        };
      }
  };
});
