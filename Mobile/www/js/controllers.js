angular.module('starter.controllers', [])

.controller('AppCtrl', function($scope, $ionicModal, $timeout) {
  // Form data for the login modal
  $scope.loginData = {};

  // Create the login modal that we will use later
  $ionicModal.fromTemplateUrl('templates/login.html', {
    scope: $scope
  }).then(function(modal) {
    $scope.modal = modal;
  });

  // Triggered in the login modal to close it
  $scope.closeLogin = function() {
    $scope.modal.hide();
  };

  // Open the login modal
  $scope.login = function() {
    $scope.modal.show();
  };

  // Perform the login action when the user submits the login form
  $scope.doLogin = function() {
    console.log('Doing login', $scope.loginData);

    // Simulate a login delay. Remove this and replace with your login
    // code if using a login system
    $timeout(function() {
      $scope.closeLogin();
    }, 1000);
  };
})

.controller('ProfilCtrl', function($scope, $http) {
    $scope.friends = [];
    $scope.profil = {};

    $scope.getFriends = function() {
      console.log("Fetching friends from server");
      $http.get('http://192.168.1.111/friends').success(function(data) {
        $scope.friends = data;
        console.log("Fetching friends succeed !");
      });
    };

    $scope.getProfil = function() {
      console.log("Fetching Profil Info")
      $http.get('http://192.168.1.111/profil').success(function(data) {
        $scope.profil = data;
        console.log("Fetching profil succeed !");
      });
    };

    $scope.getProfil();
})

.controller('FriendsCtrl', function($scope, $http) {
    $scope.friends = [];
    $scope.profil = {};

    $scope.getFriends = function() {
      console.log("Fetching friends from server");
      $http.get('http://192.168.1.111/friends').success(function(data) {
        $scope.friends = data;
        console.log("Fetching friends succeed !");
      });
    };

    $scope.getProfil = function() {
      console.log("Fetching Profil Info")
      $http.get('http://192.168.1.111/profil').success(function(data) {
        $scope.profil = data;
        console.log("Fetching profil succeed !");
      });
    }
})

.controller('MapCtrl', function($scope, $http) {
    console.log("MapCtrl Loading !");
    $scope.message = [];

    $scope.init = function() {
        $scope.getMessage();
        var myLatlng = new google.maps.LatLng(37.3000, -120.4833);
 
        var mapOptions = {
            center: myLatlng,
            zoom: 16,
            mapTypeId: google.maps.MapTypeId.ROADMAP
        };
 
        var map = new google.maps.Map(document.getElementById("map-canvas"), mapOptions);
 
        navigator.geolocation.getCurrentPosition(function(pos) {
            map.setCenter(new google.maps.LatLng(pos.coords.latitude, pos.coords.longitude));

            var markers = [];

            for (var i = 0; i < $scope.message.length; i++) {
              markers.push(new google.maps.Marker({
                  position: new google.maps.LatLng($scope.message[i].position.lat, $scope.message[i].position.long),
                  map: map,
                  title: $scope.message[i].from.username
              }));
            };
            for(var i = 0; i < markers.length; i++) {
              google.maps.event.addListener(markers[i], 'click', function(data) {
                console.log($scope.message[i-1]);
                map.setCenter(markers[i-1].getPosition());
              });
            }

        });
    };

    $scope.getMessage = function() {
      $http.get("http://192.168.1.111/message").success(function(data) {
        $scope.message = data;
      });
    };
});