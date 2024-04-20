require 'rails_helper'

RSpec.describe "Users", type: :request do
  describe "GET /new" do
    it "returns http success" do
      get "/users/new"
      expect(response).to have_http_status(:success)
    end
  end

  describe "GET /createrails" do
    it "returns http success" do
      get "/users/createrails"
      expect(response).to have_http_status(:success)
    end
  end

  describe "GET /generate" do
    it "returns http success" do
      get "/users/generate"
      expect(response).to have_http_status(:success)
    end
  end

  describe "GET /controller" do
    it "returns http success" do
      get "/users/controller"
      expect(response).to have_http_status(:success)
    end
  end

  describe "GET /Users" do
    it "returns http success" do
      get "/users/Users"
      expect(response).to have_http_status(:success)
    end
  end

  describe "GET /new" do
    it "returns http success" do
      get "/users/new"
      expect(response).to have_http_status(:success)
    end
  end

  describe "GET /create" do
    it "returns http success" do
      get "/users/create"
      expect(response).to have_http_status(:success)
    end
  end

end
