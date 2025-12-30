// 主要交互逻辑

// DOM 加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 初始化所有功能
    initNavigation();
    initCopyButtons();
    initBackToTop();
    initMobileMenu();
});

// 导航高亮和平滑滚动
function initNavigation() {
    const navLinks = document.querySelectorAll('nav a[href^="#"]');
    const sections = document.querySelectorAll('section[id]');
    
    // 点击导航链接平滑滚动
    navLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const targetId = this.getAttribute('href');
            const targetSection = document.querySelector(targetId);
            
            if (targetSection) {
                const headerOffset = 64; // header 高度
                const elementPosition = targetSection.offsetTop;
                const offsetPosition = elementPosition - headerOffset;
                
                window.scrollTo({
                    top: offsetPosition,
                    behavior: 'smooth'
                });
                
                // 移动端关闭菜单
                const nav = document.querySelector('nav');
                if (nav && nav.classList.contains('open')) {
                    nav.classList.remove('open');
                }
            }
        });
    });
    
    // 滚动时高亮当前章节
    window.addEventListener('scroll', function() {
        let current = '';
        const scrollPosition = window.pageYOffset + 100;
        
        sections.forEach(section => {
            const sectionTop = section.offsetTop;
            const sectionHeight = section.clientHeight;
            
            if (scrollPosition >= sectionTop && scrollPosition < sectionTop + sectionHeight) {
                current = section.getAttribute('id');
            }
        });
        
        navLinks.forEach(link => {
            link.classList.remove('active');
            if (link.getAttribute('href') === `#${current}`) {
                link.classList.add('active');
            }
        });
    });
}

// 代码复制功能
function initCopyButtons() {
    const codeBlocks = document.querySelectorAll('.code-block');
    
    codeBlocks.forEach(block => {
        const button = block.querySelector('.copy-btn');
        const code = block.querySelector('code');
        
        if (button && code) {
            button.addEventListener('click', function() {
                // 复制代码到剪贴板
                const text = code.textContent;
                
                navigator.clipboard.writeText(text).then(() => {
                    // 更新按钮文字
                    const originalText = button.textContent;
                    button.textContent = '已复制!';
                    button.classList.add('copied');
                    
                    // 2秒后恢复
                    setTimeout(() => {
                        button.textContent = originalText;
                        button.classList.remove('copied');
                    }, 2000);
                }).catch(err => {
                    console.error('复制失败:', err);
                    alert('复制失败，请手动复制');
                });
            });
        }
    });
}

// 返回顶部按钮
function initBackToTop() {
    const backToTopButton = document.querySelector('.back-to-top');
    
    if (backToTopButton) {
        // 滚动时显示/隐藏按钮
        window.addEventListener('scroll', function() {
            if (window.pageYOffset > 300) {
                backToTopButton.classList.add('visible');
            } else {
                backToTopButton.classList.remove('visible');
            }
        });
        
        // 点击返回顶部
        backToTopButton.addEventListener('click', function() {
            window.scrollTo({
                top: 0,
                behavior: 'smooth'
            });
        });
    }
}

// 移动端菜单
function initMobileMenu() {
    const hamburger = document.querySelector('.hamburger');
    const nav = document.querySelector('nav');
    
    if (hamburger && nav) {
        hamburger.addEventListener('click', function() {
            nav.classList.toggle('open');
            
            // 更新 aria-expanded 属性
            const isOpen = nav.classList.contains('open');
            hamburger.setAttribute('aria-expanded', isOpen);
        });
        
        // 点击页面其他地方关闭菜单
        document.addEventListener('click', function(e) {
            if (!hamburger.contains(e.target) && !nav.contains(e.target)) {
                nav.classList.remove('open');
                hamburger.setAttribute('aria-expanded', 'false');
            }
        });
    }
}

// 平台检测（可选功能 - 高亮推荐下载选项）
function detectPlatform() {
    const platform = navigator.platform.toLowerCase();
    const userAgent = navigator.userAgent.toLowerCase();
    
    if (platform.includes('win')) {
        return 'windows';
    } else if (platform.includes('mac')) {
        // 检测是否为 Apple Silicon
        if (userAgent.includes('arm')) {
            return 'macos-arm64';
        }
        return 'macos-amd64';
    } else if (platform.includes('linux')) {
        return 'linux';
    }
    
    return 'unknown';
}

// 可选：高亮推荐的下载选项
function highlightRecommendedDownload() {
    const platform = detectPlatform();
    const downloadCards = document.querySelectorAll('.download-card');
    
    downloadCards.forEach(card => {
        const cardPlatform = card.dataset.platform;
        if (cardPlatform === platform) {
            card.classList.add('recommended');
        }
    });
}
