(function() {
	var app = angular.module('lifestamp', []);

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

		/*$scope.getMessage = function() {
			if (navigator.geolocation && $scope.newTag) {
				navigator.geolocation.getCurrentPosition(sendMessage, errorPosition);
			} else {
				$scope.notification = {
					 visible: true,
					 text: "Browser not Compatible"
				};
				$scope.$apply();
			};
		};*/


		var sendMessage = function(position) {
			$scope.place = {
				From: {
					Username: $('#active').val(),
				},
				To: [{
						Username: $("#friendName").val(),
					},
				],
				Message: $("#message").val(),
				Position : {
					Longitude: position.coords.longitude,
					Latitude: position.coords.latitude
				}
			};
			$http.post('place', $scope.place).success(function() {
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
			console.log($scope.messages);
			$http.post('location', current_position).success(function(data) {
				console.log(data);
				if(data.Message) {
					$scope.messages.push(data);
					console.log(data);
				};
			});
			console.log("COORDS SENT");
		};


		$interval(function() {
			$scope.getPosition();
		}, 2500); // Refresh Time


		var errorPosition = function(error) {
			$scope.err = error;
			console.log(error);
		};



	});

})();
