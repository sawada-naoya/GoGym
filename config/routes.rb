Rails.application.routes.draw do
  root 'top#index'

  get '/login', to: 'user_sessions#new', as: :login
  post '/login', to: 'user_sessions#create'
  delete '/logout', to: 'user_sessions#destroy', as: :logout

  resources :top, only: %i[index]
  resources :users, only: %i[new create]


  resources :gyms do
    collection do
      get 'search'
    end
  end
end
