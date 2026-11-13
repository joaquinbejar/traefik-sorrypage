# Traefik SorryPage Middleware Plugin

This Traefik middleware plugin allows you to configure sorrypage responses for your routers.

SorryPage mode will be triggered if `enabled` is set to `true` and if the file configured for `triggerFilename` exists.

## Static Configuration

### FILE

```yaml
experimental:
  plugins:
    traefik-sorrypage:
      moduleName: github.com/joaquinbejar/traefik-sorrypage
      version: v0.1.0
```

### CLI

```shell
--experimental.plugins.traefik-sorrypage.modulename=github.com/joaquinbejar/traefik-sorrypage
--experimental.plugins.traefik-sorrypage.version=v0.1.0
```

## Dynamic Configuration

### FILE

```yaml
http:
  services:
    service1:
      loadBalancer:
        servers:
          - url: "http://service1:8080/"
    service2:
      loadBalancer:
        servers:
          - url: "http://service2:8081/"
    sorrypage:
      loadBalancer:
        servers:
          - url: "http://sorrypage:80/"         
            
  routers:
    service1-router:
      rule: "Host(`service1`)"
      service: "service1"
      middlewares:
        - sorrypage
    service2-router:
      rule: "Host(`service2`)"
      service: "service2"
      middlewares:
        - sorrypage
    sorrypage-router:
      rule: "Host(`sorrypage`)"
      service: "sorrypage"

  middlewares:
    sorrypage:
      plugin:
        traefik-sorrypage:
          enabled: true
          RedirectService: 'sorrypage'
```

## Contribution and Contact

### **Contact Information**
- **Author**: Joaquín Béjar García
- **Email**: jb@taunais.com
- **Telegram**: [@joaquin_bejar](https://t.me/joaquin_bejar)
- **Repository**: <https://github.com/joaquinbejar/traefik-sorrypage>

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

<!-- related-projects:start -->
## Related projects

No other repositories by the same author depend on this project, and it depends on none of them. See all of them at <https://github.com/joaquinbejar?tab=repositories>.

<!-- related-projects:end -->
