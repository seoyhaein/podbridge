package podbridge

/*
	이미지를 어떻게 구성하는가가 핵심이 될듯하다.
	buildah 를 이용해서 이미지를 코드상에서 만들어 주는 방법과 기본이미지에 추가 이미지를 넣는 방법등을 생각해보자.

	이미지 구성하는 방법은 초보적으로 파악했다.

	- baseimage 선택할 수 있도록 해야한다.

	- image 는 buildah 중심으로 만들고, command 실행은 podman 중심으로 한다.

	heartbeat 문제 해결해야함. -- shell script 를 어떻게 처리할지 고민해야한다. executor 를 넣는 방법도 생각해보겠지만, 좀 부정적이다.
*/
