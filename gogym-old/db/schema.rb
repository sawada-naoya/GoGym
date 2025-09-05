# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[7.0].define(version: 2024_09_01_120402) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "authentications", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.string "provider", null: false
    t.string "uid", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["provider", "uid"], name: "index_authentications_on_provider_and_uid", unique: true
    t.index ["user_id"], name: "index_authentications_on_user_id"
  end

  create_table "favorites", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.bigint "gym_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["gym_id"], name: "index_favorites_on_gym_id"
    t.index ["user_id", "gym_id"], name: "index_favorites_on_user_id_and_gym_id", unique: true
    t.index ["user_id"], name: "index_favorites_on_user_id"
  end

  create_table "gym_tags", force: :cascade do |t|
    t.bigint "gym_id", null: false
    t.bigint "tag_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["gym_id", "tag_id"], name: "index_gym_tags_on_gym_id_and_tag_id", unique: true
    t.index ["gym_id"], name: "index_gym_tags_on_gym_id"
    t.index ["tag_id"], name: "index_gym_tags_on_tag_id"
  end

  create_table "gyms", force: :cascade do |t|
    t.string "name", null: false
    t.string "membership_fee"
    t.string "business_hours"
    t.string "access"
    t.string "remarks"
    t.string "website"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.integer "view_count", default: 0, null: false
    t.bigint "user_id", null: false
    t.index ["user_id"], name: "index_gyms_on_user_id"
  end

  create_table "inquiries", force: :cascade do |t|
    t.string "name", null: false
    t.string "email", null: false
    t.text "inquiry_content", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "locations", force: :cascade do |t|
    t.string "address", null: false
    t.float "latitude", null: false
    t.float "longitude", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.integer "gym_id"
  end

  create_table "reviews", force: :cascade do |t|
    t.string "title"
    t.text "content"
    t.float "rating"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "gym_id", null: false
    t.bigint "user_id", null: false
    t.string "image"
    t.index ["gym_id"], name: "index_reviews_on_gym_id"
    t.index ["user_id"], name: "index_reviews_on_user_id"
  end

  create_table "tags", force: :cascade do |t|
    t.string "name", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "users", force: :cascade do |t|
    t.string "email", null: false
    t.string "crypted_password"
    t.string "salt"
    t.string "name", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["email"], name: "index_users_on_email", unique: true
  end

  add_foreign_key "authentications", "users"
  add_foreign_key "favorites", "gyms"
  add_foreign_key "favorites", "users"
  add_foreign_key "gym_tags", "gyms"
  add_foreign_key "gym_tags", "tags"
  add_foreign_key "gyms", "users"
  add_foreign_key "reviews", "gyms"
  add_foreign_key "reviews", "users"
end
