const lightCodeTheme = require('prism-react-renderer').themes.dracula;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

// With JSDoc @type annotations, IDEs can provide config autocompletion
/** @type {import('@docusaurus/types').DocusaurusConfig} */
(module.exports = {
  title: 'Shield',
  tagline: 'Identity and authorization for your APIs',
  url: 'https://goto.github.io/',
  baseUrl: '/shield/',
  onBrokenLinks: 'throw',
  markdown: { hooks: { onBrokenMarkdownLinks: 'warn' } },
  favicon: 'img/favicon.ico',
  organizationName: 'goto',
  projectName: 'shield',

  // Explicit plugins preserve the classic layout without unused SVG-to-React
  // tooling (the sole dependency path to vulnerable deepmerge).
  themes: [
    ['@docusaurus/theme-classic', {
      customCss: [
        require.resolve('./src/css/theme.css'),
        require.resolve('./src/css/custom.css'),
        require.resolve('./src/css/icons.css'),
      ],
    }],
  ],
  plugins: [
    ['@docusaurus/plugin-content-docs', {
      routeBasePath: '/',
      sidebarPath: require.resolve('./sidebars.js'),
      editUrl: 'https://github.com/goto/shield/edit/master/docs/',
      sidebarCollapsed: true,
      breadcrumbs: false,
    }],
    '@docusaurus/plugin-content-pages',
    ['@docusaurus/plugin-google-gtag', { trackingID: 'G-XXX' }],
    ...(process.env.NODE_ENV === 'production'
      ? ['@docusaurus/plugin-sitemap']
      : ['@docusaurus/plugin-debug']),
  ],

  themeConfig:
    ({
      colorMode: {
        defaultMode: 'light',
        respectPrefersColorScheme: true,
      },
      navbar: {
        title: 'Shield',
        logo: { src: 'img/shield.svg', },
        hideOnScroll: true,
        items: [
          {
            type: 'doc',
            docId: 'introduction',
            position: 'right',
            label: 'Documentation',
          },
          { to: 'support', label: 'Support', position: 'right' },
          {
            href: 'https://bit.ly/2RzPbtn',
            position: 'right',
            className: 'header-slack-link',
          },
          {
            href: 'https://github.com/goto/shield',
            className: 'navbar-item-github',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'light',
        copyright: `Gotocompany © ${new Date().getFullYear()}`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      announcementBar: {
        id: 'star-repo',
        content: '⭐️ If you like Shield, give it a star on <a target="_blank" rel="noopener noreferrer" href="https://github.com/goto/shield">GitHub</a>! ⭐',
        backgroundColor: '#222',
        textColor: '#eee',
        isCloseable: true,
      },
    }),
});
