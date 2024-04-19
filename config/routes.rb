Rails.application.routes.draw do
  root 'top#index'

  resources :top

  resources :gyms do
    collection do
      get 'search'
    end
  end
end
