var app = angular.module('app', ['ngRoute']);
var url = '/';

app.config(function($httpProvider) {
  if (window.localStorage.name && window.localStorage.token) {
    $httpProvider.defaults.headers.post.Name = window.localStorage.name;
    $httpProvider.defaults.headers.post.Token = window.localStorage.token;
  }
});

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider
    .when('/', {
      templateUrl: 'views/home.html',
      controller: 'HomeController'
    })
    .when('/t/:id', {
      templateUrl: 'views/section.html',
      controller: 'SectionController'
    })
    .when('/s/:keyword', {
      templateUrl: 'views/section.html',
      controller: 'SectionController'
    })
    .when('/login', {
      templateUrl: 'views/login.html',
      controller: 'LoginController'
    })
    .when('/a', {
      templateUrl: 'views/editor.html',
      controller: 'AddController'
    })
    .when('/e/:id', {
      templateUrl: 'views/editor.html',
      controller: 'UpdateController'
    });

}]);

app.controller('NavController', function($scope, $rootScope, $http, $window, $location, $routeParams) {
  $scope.mainput = '';
  $rootScope.auth = false;
  $rootScope.path = '/';

  $scope.search = function() {
    var input = $scope.mainput;

    if (input === '') {
      $location.path('/');
    } else {
      $location.path('/s/' + input);
    }
  };

  $scope.login = function() {
    $rootScope.path = $location.path();
    $location.path('/login');
  };

  $scope.logout = function() {
    $http.post(url + "logout.php").success(
      function(data) {
        if (data) {
          $window.localStorage.token = 0;
          $rootScope.auth = false;
        }
      }
    );
  };

  (function auth() {
    $http.post(url + 'auth.php').success(
      function(data) {
        if (data) {
          $rootScope.auth = true;
          // $scope.commit = data;
        } else {
          $rootScope.auth = false;
        }
      }
    );
  })();

});

app.controller('HomeController', function($scope, $http, $window, $location, $routeParams) {
  $scope.title = [];
  // $scope.commit = [];

  (function init() {
    $http.post(url + 'list', {
      title: 1
    }).success(
      function(data) {
        $scope.title = data;
      }
    );
  })();

});

app.controller('SectionController', function($scope, $rootScope, $http, $window, $location, $routeParams) {
  $scope.fragment = [];

  $scope.init = function(data) {
    $http.post(url + 'get', data).success(
      function(data) {
        $scope.fragment = data;
      }
    );
  };

  $scope.e = function(id) {
    $rootScope.path = $location.path();
    $location.path('/e/' + id);
  };

  if ($routeParams.id) {
    data = {
      id: $routeParams.id,
      keyword: ''
    };
    $scope.init(data);
  } else if ($routeParams.keyword) {
    $scope.mainput = $routeParams.keyword;
    data = {
      id: '',
      keyword: $routeParams.keyword
    };
    $scope.init(data);
  }

});

app.controller('LoginController', function($scope, $rootScope, $http, $window, $location) {

  $scope.login = function() {
    $http.post(url + "login.php", {
      name: $scope.name,
      passwd: $scope.passwd
    }).success(
      function(data) {
        if (data) {
          $window.localStorage.name = $scope.name;
          $window.localStorage.token = data.token;
          $http.defaults.headers.post.Name = $scope.name;
          $http.defaults.headers.post.Token = data.token;
          $rootScope.auth = true;
          // console.log($rootScope.path);
          if ($rootScope.path) {
            $location.path($rootScope.path).replace();
          } else {
            $location.path('/').replace();
          }
        }
      }
    );

    $scope.passwd = '';
  };
});

app.controller('AddController', function($scope, $http, $location) {

  $scope.title = '';
  $scope.content = ''; // or undefined

  $scope.save = function() {
    if ($scope.title !== '') {
      var data = {
        title: $scope.title,
        content: $scope.content,
        id: ''
      };
      $http.post(url + "put", data).success(
        function(data) {
          if (data) {
            $location.path('/t/' + data.id).replace();
            $scope.title = '';
            $scope.content = '';
          }
        }
      );
    }
  };

});

app.controller('UpdateController', function($scope, $rootScope, $http, $location, $routeParams) {

  $scope.id = 0;
  $scope.title = '';
  $scope.content = ''; // or undefined

  (function init() {
    $http.post(url + 'get', {
      id: $routeParams.id
    }).success(
      function(data) {
        if (data) {
          $scope.title = data[0].title;
          $scope.content = data[0].content;
          $scope.id = data[0].id;
        } else {
          $location.path($rootScope.path).replace();
        }
      }
    );
  })();

  $scope.save = function() {
    if ($scope.title !== '') {
      var data = {
        title: $scope.title,
        content: $scope.content,
        id: $scope.id
      };
      $http.post(url + "put", data).success(
        function(data) {
          if (data) {
            $location.path('/t/' + $scope.id).replace();
            $scope.id = 0;
            $scope.title = '';
            $scope.content = '';
          }
        }
      );
    }
  };

});


app.filter('markdown', function() {
  var converter = new Showdown.converter();
  return function(input) {
    return converter.makeHtml(input || '');
  };
});

// 有Bug
// app.filter('highlight', function($sce) {
//   return function(input, s) {
//     if (s.length > 1) {
//       var re = new RegExp("(" + s + ")", "gi");
//       return $sce.trustAsHtml(input.replace(re, '<span class="search">$1</span>'));
//     } else { // undefined or ''
//       return $sce.trustAsHtml(input);
//     }
//   };
// });

app.filter('trustAsHtml', function($sce) {
  return function(input) {
    return $sce.trustAsHtml(input);
  };
});
