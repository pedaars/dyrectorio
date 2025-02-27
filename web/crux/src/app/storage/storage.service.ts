import { Injectable } from '@nestjs/common'
import { Identity } from '@ory/kratos-client'
import PrismaService from 'src/services/prisma.service'
import TeamRepository from '../team/team.repository'
import { CreateStorageDto, StorageDetailsDto, StorageDto, StorageOptionDto, UpdateStorageDto } from './storage.dto'
import StorageMapper from './storage.mapper'

@Injectable()
export default class StorageService {
  constructor(
    private teamRepository: TeamRepository,
    private prisma: PrismaService,
    private mapper: StorageMapper,
  ) {}

  async getStorages(teamSlug: string): Promise<StorageDto[]> {
    const storages = await this.prisma.storage.findMany({
      where: {
        team: {
          slug: teamSlug,
        },
      },
    })

    return storages.map(it => this.mapper.listItemToDto(it))
  }

  async getStorageDetails(id: string): Promise<StorageDetailsDto> {
    const storage = await this.prisma.storage.findUniqueOrThrow({
      include: {
        _count: {
          select: {
            containerConfigs: true,
            instanceConfigs: true,
          },
        },
      },
      where: {
        id,
      },
    })

    return this.mapper.detailsToDto(storage)
  }

  async createStorage(teamSlug: string, req: CreateStorageDto, identity: Identity): Promise<StorageDetailsDto> {
    const teamId = await this.teamRepository.getTeamIdBySlug(teamSlug)

    const storage = await this.prisma.storage.create({
      data: {
        name: req.name,
        description: req.description,
        icon: req.icon ?? null,
        url: req.url,
        accessKey: req.accessKey,
        secretKey: req.secretKey,
        teamId,
        createdBy: identity.id,
      },
    })

    return this.mapper.detailsToDto({
      ...storage,
      _count: {
        containerConfigs: 0,
        instanceConfigs: 0,
      },
    })
  }

  async updateStorage(id: string, req: UpdateStorageDto, identity: Identity): Promise<void> {
    await this.prisma.storage.update({
      where: {
        id,
      },
      data: {
        name: req.name,
        description: req.description,
        icon: req.icon ?? null,
        url: req.url,
        accessKey: req.accessKey,
        secretKey: req.secretKey,
        updatedBy: identity.id,
      },
    })
  }

  async deleteStorage(id: string): Promise<void> {
    await this.prisma.storage.delete({
      where: {
        id,
      },
    })
  }

  async getStorageOptions(teamSlug: string): Promise<StorageOptionDto[]> {
    const storages = await this.prisma.storage.findMany({
      where: {
        team: {
          slug: teamSlug,
        },
      },
      select: {
        id: true,
        name: true,
      },
    })

    return storages
  }
}
