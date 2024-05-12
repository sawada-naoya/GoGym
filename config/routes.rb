Rails.application.routes.draw do
  root 'top#index'

  get '/login', to: 'user_sessions#new', as: :login
  post '/login', to: 'user_sessions#create'
  delete '/logout', to: 'user_sessions#destroy', as: :logout

  resources :top, only: %i[index]
  resources :users do
    get 'reviews', to: 'reviews#user_index', as: 'user_reviews'
  end
  resources :locations


  resources :gyms do
    member do
      get 'images', to: 'gyms#images_index', as: 'images'
    end
    collection do
      get 'search'
    end
    resources :reviews, shallow: true
  end
end
