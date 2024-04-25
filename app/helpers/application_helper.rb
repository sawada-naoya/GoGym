module ApplicationHelper
  def page_title(title = '')
    base_title = 'GoGym'
    title.present? ? "#{title} | #{base_title}" : base_title
  end

  def active_if(path)
    path == controller_path ? 'active' : ''
  end
end
