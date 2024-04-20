Rails.application.routes.draw do
  root 'top#index'

  resources :top, only: %i[index]
  resources :users, only: %i[new create]

  resources :gyms do
    collection do
      get 'search'
    end
  end
end
