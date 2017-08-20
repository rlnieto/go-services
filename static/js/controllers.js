// Registering the parkingCtrl to the parking module
parking.controller("parkingCtrl", function ($scope, $filter) {

  // Binding the car's array to the scope
  $scope.cars = [];
/*    $scope.cars = [
    {plate: '6MBV006', color: 'White'},
    {plate: '5BBM299', color: 'Red'},
    {plate: '5AOJ230', color: 'Black'}
  ];
*/
  $scope.appTitle = $filter("uppercase")("[Packt] Parking");
  $scope.colors = ["White", "Black", "Blue", "Red", "Silver"];
  $scope.alertTopic = "Something went wrong";
  $scope.descriptionTopic = "You must inform de plate and the color of the car";
  $scope.showAlert=true;

  $scope.closeAlert = function(){
    $scope.showAlert=false;
  };

  // Binding the park function to the scope
  $scope.park = function (car) {
    car.entrance = new Date();
    $scope.cars.push(angular.copy(car));
    delete $scope.car;
    $scope.carForm.$setPristine();
  };
});
