Rails.application.routes.draw do
  get 'oauths/oauth'
  get 'oauths/callback'
  mount LetterOpenerWeb::Engine, at: '/letter_opener' if Rails.env.development?
  require 'sidekiq/web'
  mount Sidekiq::Web => '/sidekiq'

  root 'top#index'

  get '/login', to: 'user_sessions#new', as: :login
  post '/login', to: 'user_sessions#create'
  delete '/logout', to: 'user_sessions#destroy', as: :logout
  get '/terms', to: 'static_pages#terms'
  get '/privacy', to: 'static_pages#privacy'

  post "oauth/callback", to: "oauths#callback"
  get "oauth/callback", to: "oauths#callback"
  get "oauth/:provider", to: "oauths#oauth", as: :auth_at_provider

  resources :top, only: %i[index]
  resources :inquiries, only: [:new, :create]
  resources :users do
    get 'reviews', to: 'reviews#user_index', as: 'user_reviews'
    member do
      get :favorites
    end
  end
  resources :locations

  resources :tags, only: [:show] do
    get 'gyms', to: 'gyms#index', as: 'gyms_by_tag'
  end

  resources :gyms do
    member do
      get 'images', to: 'gyms#images_index', as: 'images'
    end
    collection do
      get 'search'
      get 'autocomplete'
    end
    resources :reviews, shallow: true
    resource :favorite, only: [:create, :destroy]
  end
end
