Rails.application.routes.draw do
  root 'top#index'

  resources :top
  resources :users

  resources :gyms do
    collection do
      get 'search'
    end
  end
end
