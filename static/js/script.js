(function() {
	var app = angular.module('stamp', []);

	app.controller('BuildPage', function($scope, $http, $interval) {
		$scope.place;
		$scope.showForm = false;
		$scope.newTag = false;
		$scope.notification = {
			visible: false,
			text: ""
		};
		$scope.messages = [];

		$scope.activate = function() {
			if ($scope.showForm) {
				$scope.showForm = false;
			} else {
				$scope.showForm = true;
				//$(".newTag").removeClass("animated fadeOutDown");
				$(".newTag").addClass("animated fadeInDown");
			};
		};

		$scope.getMessage = function() {
			$scope.newTag = true;
			$scope.getPosition();
		};

		$scope.getPosition = function() {
			if(navigator.geolocation && $scope.newTag == false) {
				navigator.geolocation.getCurrentPosition(sendPosition, errorPosition);
			}else if (navigator.geolocation && $scope.newTag == true) {
				navigator.geolocation.getCurrentPosition(sendMessage, errorPosition);
				$scope.newTag = false;
			} else {
				return "Browser not supported geolocation";
			};
		};

		$scope.saveMessage = function(message) {
			console.log(message);
			$http.post('save', message).success(function() {
				$.notify("Message save Successfuly !!");
			});
		};

		$scope.addFriend = function() {
			var friendName = $("#addFriend").val();
			$http.post("addfriend", friendName).success(function() {
				$.notify(friendName + " added !!")
			});
		};

		var sendMessage = function(position) {
			$scope.place = {
				From: {
					Username: $('#active').val(),
				},
				To: $("#friendName").val(),
				Message: $("#message").val(),
				Position : {
					Longitude: position.coords.longitude,
					Latitude: position.coords.latitude
				},
				Picture: "empty",
			};
			console.log($scope.place.Picture);
			$http.post('insert', $scope.place).success(function() {
				$.notify("Message Sent Successfuly !!")
				$("#friendName").val("");
				$("#message").val("");
				setTimeout(function() {
					$scope.showForm = false;
					$scope.$apply();
				}, 3000);
			});
		};


		var sendPosition = function(position) {
			var current_position = position.coords;
			//console.log($scope.messages);
			$http.post('location', current_position).success(function(data) {
				if(data.Message) {
					$scope.messages.push(data);
					console.log(data);
				};
			});
			//console.log("COORDS SENT");
		};
/*
		var findIndex = function(value, array) {
			for(var i = 0; i < array.length; i++) {
				console.log(array[i].Message);
				if (array[i].Message === value) {
					console.log(array[i].Message);
					return array[i];
				}else {
					return null;
					console.log("VALUE NOT FIND");
				};
			};
		};
*/
		$interval(function() {
			$scope.getPosition();
		}, 2500); // Refresh Time


		var errorPosition = function(error) {
			$scope.err = error;
			console.log(error);
		};



	});
	
	app.controller('Profil', function($scope, $http) {

	});
})();
