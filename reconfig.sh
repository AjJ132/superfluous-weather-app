#first start minikube
minikube start 
#ECHO WIP and EXIT
echo "WIP"
exit 0

#navigate to k8s dir and remove all deployments and services
cd k8s

kubectl delete -f .

#rebuild all docker images using buildx and push to docker hub

#Backend
cd ../Backend-For-Frontend
docker buildx build -t ajj132/super-weather-app-bff:latest . 
docker push ajj132/super-weather-app-bff:latest

#Frontend
cd ../frontend/superfluous-weather-app
npm run build
docker buildx build -t ajj132/super-weather-app-frontend:latest .
docker push ajj132/super-weather-app-frontend:latest

#login system
#login database
cd ../../login-system/user-database
docker buildx build -t ajj132/super-weather-app-login-database:latest .
docker push ajj132/super-weather-app-login-database:latest

#user manager
cd ../user-manager
docker buildx build -t ajj132/super-weather-app-login:latest .
docker push ajj132/super-weather-app-login:latest

#apply all k8s files
cd ../../k8s
kubectl apply -f .

#print done
echo "done"
