#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
HTML文档生成器配置文件
定义所有样式变量和尺寸参数，方便统一修改
"""

# 基础配置
BASE_CONFIG = {
    'font_size': '18px',
    'line_height': '1.1',
    'font_family': "-apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif",
    'background_color': '#f8f9fa',
    'text_color': '#333',
}

# 侧边栏配置
SIDEBAR_CONFIG = {
    'width': '280px',
    'background': '#fff',
    'border_color': '#e1e4e8',
    'padding': '16px',
    'shadow': '2px 0 8px rgba(0,0,0,0.08)',
    'title_font_size': '18px',
    'title_margin_bottom': '12px',
}

# 目录配置
TOC_CONFIG = {
    'item_margin': '1px 0',
    'toggle_size': '14px',
    'toggle_font_size': '11px',
    'toggle_margin_right': '3px',
    'link_padding': '3px 6px',
    'link_border_radius': '3px',
    'link_font_size': '14px',
    'level1_font_size': '15px',
    'level2_font_size': '14px',
    'level3_font_size': '13px',
    'children_margin_left': '16px',
    'children_padding_left': '6px',
}

# 主内容区配置
MAIN_CONFIG = {
    'margin_left': '280px',  # 与侧边栏宽度一致
    'padding': '8px',
    'mobile_padding': '6px',
}

# 标题配置
HEADER_CONFIG = {
    'padding': '12px',
    'border_radius': '6px',
    'margin_bottom': '10px',
    'shadow': '0 3px 12px rgba(0,0,0,0.1)',
    'title_font_size': '2.2em',
    'title_margin_bottom': '6px',
}

# 章节配置
SECTION_CONFIG = {
    'padding': '10px',
    'border_radius': '4px',
    'margin_bottom': '10px',
    'shadow': '0 2px 8px rgba(0,0,0,0.05)',
    'title_font_size': '1.6em',
    'title_margin_bottom': '8px',
    'title_border_width': '2px',
    'title_padding_bottom': '4px',
}

# 描述项配置
DESCRIPTION_CONFIG = {
    'item_margin_bottom': '10px',
    'item_padding': '8px',
    'item_border_radius': '6px',
    'item_border_width': '3px',
    'subtitle_margin_bottom': '8px',
    'subtitle_font_size': '1em',
    'text_font_size': '15px',
    'text_line_height': '1.6',
}

# 版本卡片配置
VERSION_CONFIG = {
    'card_margin_bottom': '18px',
    'card_border_radius': '8px',
    'card_shadow': '0 2px 6px rgba(0,0,0,0.05)',
    'title_padding': '14px 18px',
    'title_font_size': '1.2em',
    'description_padding': '14px 18px',
    'spec_padding': '16px 18px',
    'spec_title_margin_bottom': '16px',
    'spec_title_font_size': '1.1em',
}

# 属性概览配置
OVERVIEW_CONFIG = {
    'margin_bottom': '24px',
    'title_margin_bottom': '12px',
    'title_font_size': '1em',
    'grid_min_width': '280px',
    'grid_gap': '12px',
    'item_border_radius': '6px',
    'item_padding': '12px',
    'header_margin_bottom': '6px',
    'name_font_size': '14px',
    'desc_font_size': '13px',
    'hover_transform': 'translateY(-1px)',
    'hover_shadow': '0 3px 10px rgba(0,0,0,0.1)',
}

# 属性详情配置
PROPERTY_CONFIG = {
    'detail_margin_bottom': '16px',
    'detail_border_radius': '6px',
    'detail_shadow': '0 1px 3px rgba(0,0,0,0.05)',
    'header_padding': '6px 10px',
    'name_font_size': '1em',
    'description_padding': '6px 10px',
    'desc_item_margin_bottom': '6px',
    'desc_item_font_size': '14px',
    'desc_item_line_height': '1.5',
}

# 类型标签配置
TYPE_BADGE_CONFIG = {
    'padding': '3px 8px',
    'border_radius': '12px',
    'font_size': '13px',
    'letter_spacing': '0.3px',
}

# Release Stage标签配置
RELEASE_STAGE_CONFIG = {
    'padding': '3px 8px',
    'border_radius': '12px',
    'font_size': '12px',
    'letter_spacing': '0.3px',
    'margin_left': '6px',
}

# 约束配置
CONSTRAINT_CONFIG = {
    'section_padding': '6px 10px',
    'title_margin_bottom': '8px',
    'title_font_size': '14px',
    'list_gap': '4px',
    'item_padding': '4px 6px',
    'item_border_radius': '2px',
    'item_font_size': '14px',
    'label_margin_right': '2px',
    'value_padding': '1px 2px',
    'value_border_radius': '2px',
    'value_font_size': '13px',
    'enum_gap': '3px',
    'enum_margin_top': '2px',
    'enum_padding': '2px 6px',
    'enum_border_radius': '4px',
    'enum_font_size': '13px',
}

# 嵌套属性配置
NESTED_CONFIG = {
    'padding': '6px 8px',
    'title_margin_bottom': '6px',
    'title_font_size': '13px',
}

# 滚动条配置
SCROLLBAR_CONFIG = {
    'width': '6px',
    'track_color': '#f1f1f1',
    'thumb_color': '#c1c1c1',
    'thumb_hover_color': '#a8a8a8',
    'thumb_border_radius': '3px',
}

# 颜色配置
COLORS = {
    'primary': '#0366d6',
    'secondary': '#586069',
    'success': '#28a745',
    'danger': '#dc3545',
    'warning': '#ffc107',
    'info': '#17a2b8',
    'light': '#f8f9fa',
    'dark': '#343a40',
    'border': '#e1e4e8',
    'hover_bg': '#f1f8ff',
    'gradient_primary': 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    'gradient_secondary': 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    'blue_border': '#3498db',
    'red_border': '#e74c3c',
}

# 类型颜色配置
TYPE_COLORS = {
    'string': {'bg': '#e3f2fd', 'color': '#1976d2'},
    'number': {'bg': '#f3e5f5', 'color': '#7b1fa2'},
    'boolean': {'bg': '#e8f5e8', 'color': '#388e3c'},
    'object': {'bg': '#fff3e0', 'color': '#f57c00'},
    'array': {'bg': '#fce4ec', 'color': '#c2185b'},
    'map': {'bg': '#f1f8e9', 'color': '#689f38'},
    'enum': {'bg': '#e0f2f1', 'color': '#00695c'},
    'time': {'bg': '#e8eaf6', 'color': '#3f51b5'},
    'default': {'bg': '#f5f5f5', 'color': '#616161'},
}

# Release Stage颜色配置
RELEASE_STAGE_COLORS = {
    'experimental': {'bg': '#fff3cd', 'color': '#856404', 'border': '#ffeaa7'},
    'preview': {'bg': '#d1ecf1', 'color': '#0c5460', 'border': '#bee5eb'},
    'beta': {'bg': '#d4edda', 'color': '#155724', 'border': '#c3e6cb'},
    'ga': {'bg': '#e2e3e5', 'color': '#383d41', 'border': '#d6d8db'},
    'deprecated': {'bg': '#f8d7da', 'color': '#721c24', 'border': '#f5c6cb'},
    'default': {'bg': '#f5f5f5', 'color': '#616161', 'border': '#e0e0e0'},
}

# 响应式断点
BREAKPOINTS = {
    'mobile': '768px',
}

# 动画配置
ANIMATION_CONFIG = {
    'transition_duration': '0.2s',
    'hover_transition': '0.3s',
    'scroll_behavior': 'smooth',
} 

