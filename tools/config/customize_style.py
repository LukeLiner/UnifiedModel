#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
样式自定义示例脚本
演示如何快速修改HTML文档的样式配置
"""

import os
import shutil
from config import *

def create_custom_config(theme_name: str):
    """创建自定义配置文件"""
    
    themes = {
        'large': {
            'description': '大字体宽松布局主题',
            'changes': {
                'BASE_CONFIG': {
                    'font_size': '16px',
                },
                'SIDEBAR_CONFIG': {
                    'width': '320px',
                    'padding': '20px',
                    'title_font_size': '17px',
                },
                'MAIN_CONFIG': {
                    'margin_left': '320px',
                    'padding': '32px',
                },
                'TOC_CONFIG': {
                    'link_font_size': '15px',
                    'level1_font_size': '16px',
                    'level2_font_size': '15px',
                    'level3_font_size': '14px',
                },
                'SECTION_CONFIG': {
                    'padding': '24px',
                    'margin_bottom': '24px',
                },
            }
        },
        
        'compact': {
            'description': '紧凑布局主题',
            'changes': {
                'BASE_CONFIG': {
                    'font_size': '13px',
                },
                'SIDEBAR_CONFIG': {
                    'width': '240px',
                    'padding': '12px',
                    'title_font_size': '14px',
                },
                'MAIN_CONFIG': {
                    'margin_left': '240px',
                    'padding': '16px',
                },
                'TOC_CONFIG': {
                    'link_font_size': '12px',
                    'level1_font_size': '13px',
                    'level2_font_size': '12px',
                    'level3_font_size': '11px',
                },
                'SECTION_CONFIG': {
                    'padding': '16px',
                    'margin_bottom': '16px',
                },
                'PROPERTY_CONFIG': {
                    'detail_margin_bottom': '12px',
                },
            }
        },
        
        'dark': {
            'description': '深色主题',
            'changes': {
                'BASE_CONFIG': {
                    'background_color': '#1a1a1a',
                    'text_color': '#e0e0e0',
                },
                'SIDEBAR_CONFIG': {
                    'background': '#2d2d2d',
                    'border_color': '#404040',
                },
                'COLORS': {
                    'border': '#404040',
                    'light': '#2d2d2d',
                    'hover_bg': '#404040',
                    'gradient_primary': 'linear-gradient(135deg, #4a90e2 0%, #7b68ee 100%)',
                    'gradient_secondary': 'linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%)',
                },
            }
        },
        
        'blue': {
            'description': '蓝色主题',
            'changes': {
                'COLORS': {
                    'primary': '#007bff',
                    'gradient_primary': 'linear-gradient(135deg, #007bff 0%, #6610f2 100%)',
                    'gradient_secondary': 'linear-gradient(135deg, #17a2b8 0%, #007bff 100%)',
                    'blue_border': '#007bff',
                    'red_border': '#17a2b8',
                },
            }
        },
        
        'green': {
            'description': '绿色主题',
            'changes': {
                'COLORS': {
                    'primary': '#28a745',
                    'gradient_primary': 'linear-gradient(135deg, #28a745 0%, #20c997 100%)',
                    'gradient_secondary': 'linear-gradient(135deg, #20c997 0%, #17a2b8 100%)',
                    'blue_border': '#28a745',
                    'red_border': '#20c997',
                },
            }
        }
    }
    
    if theme_name not in themes:
        print(f"❌ 未知主题: {theme_name}")
        print(f"可用主题: {', '.join(themes.keys())}")
        return False
    
    theme = themes[theme_name]
    config_content = f'''#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
自定义配置文件 - {theme['description']}
基于 config.py 生成，可以直接修改此文件来调整样式
"""

# 导入基础配置
from config import *

# 主题: {theme_name} - {theme['description']}
'''
    
    # 生成配置覆盖代码
    for config_name, changes in theme['changes'].items():
        config_content += f'\n# 覆盖 {config_name}\n'
        for key, value in changes.items():
            if isinstance(value, str):
                config_content += f"{config_name}['{key}'] = '{value}'\n"
            else:
                config_content += f"{config_name}['{key}'] = {value}\n"
    
    # 写入自定义配置文件
    custom_config_file = f'config_{theme_name}.py'
    with open(custom_config_file, 'w', encoding='utf-8') as f:
        f.write(config_content)
    
    print(f"✅ 已创建自定义配置: {custom_config_file}")
    print(f"📝 主题描述: {theme['description']}")
    return True

def apply_custom_config(config_file: str):
    """应用自定义配置"""
    if not os.path.exists(config_file):
        print(f"❌ 配置文件不存在: {config_file}")
        return False
    
    # 备份原配置
    if os.path.exists('config.py'):
        shutil.copy('config.py', 'config_backup.py')
        print("📦 已备份原配置为: config_backup.py")
    
    # 应用新配置
    shutil.copy(config_file, 'config.py')
    print(f"✅ 已应用配置: {config_file}")
    return True

def restore_config():
    """恢复原始配置"""
    if os.path.exists('config_backup.py'):
        shutil.copy('config_backup.py', 'config.py')
        print("✅ 已恢复原始配置")
        return True
    else:
        print("❌ 未找到备份配置文件")
        return False

def demo_themes():
    """演示所有主题"""
    import subprocess
    
    themes = ['large', 'compact', 'dark', 'blue', 'green']
    input_file = 'expanded_schemas/entity_set.expanded.yaml'
    
    if not os.path.exists(input_file):
        print(f"❌ 测试文件不存在: {input_file}")
        return
    
    print("🎨 开始生成主题演示...")
    
    for theme in themes:
        print(f"\n🔄 生成 {theme} 主题...")
        
        # 创建自定义配置
        if create_custom_config(theme):
            # 应用配置
            apply_custom_config(f'config_{theme}.py')
            
            # 生成HTML
            output_file = f'demo_{theme}_theme.html'
            try:
                result = subprocess.run([
                    'python3', 'yaml_to_html_v2.py', 
                    input_file, '-o', output_file
                ], capture_output=True, text=True)
                
                if result.returncode == 0:
                    print(f"✅ 已生成: {output_file}")
                else:
                    print(f"❌ 生成失败: {result.stderr}")
            except Exception as e:
                print(f"❌ 执行失败: {e}")
    
    # 恢复原配置
    restore_config()
    
    print("\n🎉 主题演示完成！")
    print("📁 生成的文件:")
    for theme in themes:
        output_file = f'demo_{theme}_theme.html'
        if os.path.exists(output_file):
            print(f"   - {output_file}")

def main():
    import argparse
    
    parser = argparse.ArgumentParser(description='HTML文档样式自定义工具')
    parser.add_argument('action', choices=['create', 'apply', 'restore', 'demo'], 
                       help='操作类型')
    parser.add_argument('--theme', help='主题名称 (large/compact/dark/blue/green)')
    parser.add_argument('--config', help='配置文件路径')
    
    args = parser.parse_args()
    
    if args.action == 'create':
        if not args.theme:
            print("❌ 请指定主题名称: --theme <theme_name>")
            print("可用主题: large, compact, dark, blue, green")
            return
        create_custom_config(args.theme)
    
    elif args.action == 'apply':
        if not args.config:
            print("❌ 请指定配置文件: --config <config_file>")
            return
        apply_custom_config(args.config)
    
    elif args.action == 'restore':
        restore_config()
    
    elif args.action == 'demo':
        demo_themes()

if __name__ == "__main__":
    main() 